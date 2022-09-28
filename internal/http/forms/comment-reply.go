package forms

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

type CreateCommentReplyForm struct {
	ParentCommentUid string `json:"commentUid" binding:"required,len=24"`
	Email            string `json:"email" binding:"required,email"`
	Name             string `json:"name" binding:"required,max=50"`
	Content          string `json:"content" binding:"required,max=255"`

	realParentCommentUid primitive.ObjectID
	realPostUid          primitive.ObjectID
	realPostAuthorUid    primitive.ObjectID
}

func (form *CreateCommentReplyForm) Validate(
	svc *service.Service,
	ctx context.Context,
) (parentComment *models.CommentModel, err error) {
	var parentCommnetUid primitive.ObjectID

	if parentCommnetUid, err = primitive.ObjectIDFromHex(form.ParentCommentUid); err != nil {
		return nil, errors.New("invalid parent comment uid format")
	}
	if parentComment, err = findCommentForReply(svc, ctx, parentCommnetUid); err != nil {
		return nil, err
	}
	if _, err = findPostForComment(svc, ctx, parentComment.UID); err != nil {
		return nil, err
	}
	form.realParentCommentUid = parentComment.UID
	form.realPostUid = parentComment.PostUid
	form.realPostAuthorUid = parentComment.PostAuthorUid

	return parentComment, nil
}

func (form *CreateCommentReplyForm) ToCommentReplyModel() (model *models.CommentModel, err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	if len(form.realParentCommentUid) == 0 {
		return nil, errors.New("validate the form first")
	}

	return &models.CommentModel{
		UID:              primitive.NewObjectID(),
		PostUid:          form.realPostUid,
		PostAuthorUid:    form.realPostAuthorUid,
		ParentCommentUid: form.realParentCommentUid,
		Email:            form.Email,
		Name:             form.Name,
		Content:          form.Content,
		ReplyCount:       0,
		CreatedAt:        now,
		DeletedAt:        nil,
	}, nil
}

func findCommentForReply(
	svc *service.Service,
	ctx context.Context,
	formCommentUid primitive.ObjectID,
) (comment *models.CommentModel, err error) {
	if comment, err = svc.Comment.GetOne(ctx, bson.M{
		"$and": []bson.M{
			{"deletedat": bson.M{"$eq": primitive.Null{}}},
			{"_id": bson.M{"$eq": formCommentUid}}}},
	); err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, errors.New("parent comment not found")
	}

	return comment, nil
}

package forms

import (
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
	postService *service.PostService,
	commentService *service.CommentService,
) (parentComment *models.CommentModel, err error) {
	var parentCommnetUid primitive.ObjectID

	if parentCommnetUid, err = primitive.ObjectIDFromHex(form.ParentCommentUid); err != nil {
		return nil, errors.New("invalid parent comment uid format")
	}
	if parentComment, err = findCommentForReply(commentService, parentCommnetUid); err != nil {
		return nil, err
	}
	if _, err = findPostForComment(postService, parentComment.UID); err != nil {
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
	commentService *service.CommentService,
	formCommentUid primitive.ObjectID,
) (comment *models.CommentModel, err error) {
	if comment, err = commentService.GetComment(bson.M{
		"$and": []bson.M{
			{"deletedat": bson.M{"$eq": primitive.Null{}}},
			{"_id": bson.M{"$eq": formCommentUid}}},
	}); err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, errors.New("parent comment not found")
	}

	return comment, nil
}

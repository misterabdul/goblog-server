package forms

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

type ReplyCommmentForm struct {
	CommentUid string `json:"commentUid" binding:"required,len=24"`
	Email      string `json:"email" binding:"required,email"`
	Name       string `json:"name" binding:"required,max=50"`
	Content    string `json:"content" binding:"required,max=255"`

	parentComment *models.CommentModel
}

func (form *ReplyCommmentForm) Validate(
	commentService *service.Service,
) (err error) {
	if form.parentComment, err = findComment(commentService, form.CommentUid); err != nil {
		return err
	}

	return nil
}

func (form *ReplyCommmentForm) ToCommentModel() (model *models.CommentModel, err error) {
	var (
		now           = primitive.NewDateTimeFromTime(time.Now())
		parentComment = form.parentComment
		replies       = parentComment.Replies
	)

	if form.parentComment == nil {
		return nil, errors.New("validate the form first")
	}

	replies = append(replies, models.CommentReplyModel{
		Email:     form.Email,
		Name:      form.Name,
		Content:   form.Content,
		CreatedAt: now,
		DeletedAt: nil})
	parentComment.Replies = replies

	return parentComment, nil
}

func findComment(
	commentService *service.Service,
	formCommentUid string,
) (comment *models.CommentModel, err error) {
	if comment, err = commentService.GetComment(bson.M{
		"$and": []bson.M{
			{"deletedat": bson.M{"$eq": primitive.Null{}}},
			{"$or": []bson.M{
				{"_id": bson.M{"$eq": formCommentUid}}}}},
	}); err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, errors.New("comment not found")
	}

	return comment, nil
}

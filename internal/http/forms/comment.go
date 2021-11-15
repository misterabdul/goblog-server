package forms

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/models"
)

type CreateCommentForm struct {
	PostUid string `json:"postUid" binding:"required,printascii,max=24"`
	Email   string `json:"email" binding:"required,email"`
	Name    string `json:"name" binding:"required,printascii,max=50"`
	Content string `json:"content" bindinng:"required,printascii,max=255"`
}

type ReplyCommmentForm struct {
	CommentUid string `json:"commentUid" binding:"required,printascii,max=24"`
	Email      string `json:"email" binding:"required,email"`
	Name       string `json:"name" binding:"required,printascii,max=50"`
	Content    string `json:"content" binding:"required,printascii,max=255"`
}

func CreateCommentModel(form *CreateCommentForm) (model *models.CommentModel, err error) {
	var (
		postUid primitive.ObjectID
		now     primitive.DateTime = primitive.NewDateTimeFromTime(time.Now())
	)

	if postUid, err = primitive.ObjectIDFromHex(form.PostUid); err != nil {
		return nil, err
	}

	return &models.CommentModel{
		UID:       primitive.NewObjectID(),
		PostUid:   postUid,
		Email:     form.Email,
		Name:      form.Name,
		Content:   form.Content,
		Replies:   []models.CommentReplyModel{},
		CreatedAt: now,
		DeletedAt: nil}, nil
}

func ReplyCommentModel(
	form *ReplyCommmentForm,
	comment *models.CommentModel,
) (model *models.CommentModel) {
	now := primitive.NewDateTimeFromTime(time.Now())
	replies := comment.Replies
	replies = append(replies, models.CommentReplyModel{
		Email:     form.Email,
		Name:      form.Name,
		Content:   form.Content,
		CreatedAt: now,
		DeletedAt: nil})
	comment.Replies = replies

	return comment
}

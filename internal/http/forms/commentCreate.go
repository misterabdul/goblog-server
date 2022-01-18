package forms

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

type CreateCommentForm struct {
	PostUid string `json:"postUid" binding:"required,alphanum,len=24"`
	Email   string `json:"email" binding:"required,email"`
	Name    string `json:"name" binding:"required,max=50"`
	Content string `json:"content" bindinng:"required,max=255"`

	realPostUid primitive.ObjectID
}

func (form *CreateCommentForm) Validate(
	commentService *service.Service,
) (err error) {
	var (
		post *models.PostModel
	)

	if post, err = findPostForComment(commentService, form.PostUid); err != nil {
		return err
	}
	form.realPostUid = post.UID

	return nil
}

func (form *CreateCommentForm) ToCommentModel() (model *models.CommentModel, err error) {
	var (
		now = primitive.NewDateTimeFromTime(time.Now())
	)

	if len(form.realPostUid) == 0 {
		return nil, errors.New("validate the form first")
	}

	return &models.CommentModel{
		UID:       primitive.NewObjectID(),
		PostUid:   form.realPostUid,
		Email:     form.Email,
		Name:      form.Name,
		Content:   form.Content,
		Replies:   []models.CommentReplyModel{},
		CreatedAt: now,
		DeletedAt: nil}, nil
}

func findPostForComment(
	commentService *service.Service,
	formPostUid string,
) (post *models.PostModel, err error) {
	if post, err = commentService.GetPost(bson.M{
		"$and": []bson.M{
			{"deletedat": bson.M{"$eq": primitive.Null{}}},
			{"publishedat": bson.M{"$ne": primitive.Null{}}},
			{"$or": []bson.M{
				{"_id": bson.M{"$eq": formPostUid}},
				{"slug": bson.M{"$eq": formPostUid}}}}},
	}); err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("post not found")
	}

	return post, nil
}

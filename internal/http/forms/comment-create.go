package forms

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

type CreateCommentForm struct {
	PostUid string `json:"postUid" binding:"required,alphanum,len=24"`
	Email   string `json:"email" binding:"required,email"`
	Name    string `json:"name" binding:"required,max=50"`
	Content string `json:"content" bindinng:"required,max=255"`

	realPostUid       primitive.ObjectID
	realPostAuthorUid primitive.ObjectID
}

func (form *CreateCommentForm) Validate(
	postService *service.PostService,
) (post *models.PostModel, err error) {
	var postUid primitive.ObjectID

	if postUid, err = primitive.ObjectIDFromHex(form.PostUid); err != nil {
		return nil, errors.New("invalid post uid format")
	}
	if post, err = findPostForComment(postService, postUid); err != nil {
		return nil, err
	}
	form.realPostUid = post.UID
	form.realPostAuthorUid = post.Author.UID

	return post, nil
}

func (form *CreateCommentForm) ToCommentModel() (model *models.CommentModel, err error) {
	var now = primitive.NewDateTimeFromTime(time.Now())

	if len(form.realPostUid) == 0 {
		return nil, errors.New("validate the form first")
	}
	if len(form.realPostAuthorUid) == 0 {
		return nil, errors.New("validate the form first")
	}

	return &models.CommentModel{
		UID:              primitive.NewObjectID(),
		PostUid:          form.realPostUid,
		PostAuthorUid:    form.realPostAuthorUid,
		ParentCommentUid: nil,
		Email:            form.Email,
		Name:             form.Name,
		Content:          form.Content,
		ReplyCount:       0,
		CreatedAt:        now,
		DeletedAt:        nil,
	}, nil
}

func findPostForComment(
	postService *service.PostService,
	formPostUid primitive.ObjectID,
) (post *models.PostModel, err error) {
	if post, err = postService.GetPost(bson.M{
		"$and": []bson.M{
			{"deletedat": bson.M{"$eq": primitive.Null{}}},
			{"publishedat": bson.M{"$ne": primitive.Null{}}},
			{"_id": bson.M{"$eq": formPostUid}}},
	}); err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("post not found")
	}

	return post, nil
}

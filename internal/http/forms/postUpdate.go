package forms

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/models"
)

type UpdatePostForm struct {
	Slug               string   `json:"slug" binding:"omitempty,alphanum,max=100"`
	Title              string   `json:"title" binding:"omitempty,max=100"`
	Description        string   `json:"description" binding:"omitempty,max=255"`
	FeaturingImagePath string   `json:"featuringImagePath" binding:"omitempty,url"`
	Categories         []string `json:"categories" binding:"omitempty,dive,len=12"`
	Tags               []string `json:"tags" binding:"omitempty,dive,max=32"`
	Content            string   `json:"content" binding:"omitempty"`
	PublishNow         bool     `json:"publishNow" binding:"omitempty,boolean"`

	realCategories []*models.CategoryModel
}

func (form *UpdatePostForm) Validate(
	ctx context.Context,
	dbConn *mongo.Database,
) (err error) {
	if err = checkPostSlug(ctx, dbConn, form.Slug); err != nil {
		return err
	}
	if form.realCategories, err = findCategories(ctx, dbConn, form.Categories); err != nil {
		return err
	}

	return nil
}

func (form *UpdatePostForm) ToPostModel(
	post *models.PostModel,
	postContent *models.PostContentModel,
) (
	updatedPost *models.PostModel,
	updatedPostContent *models.PostContentModel,
	err error,
) {
	var (
		categories = []models.CategoryCommonModel{}
		now        = primitive.NewDateTimeFromTime(time.Now())
	)

	if len(form.Slug) > 0 {
		post.Slug = form.Slug
	}
	if len(form.Title) > 0 {
		post.Title = form.Title
	}
	if len(form.Description) > 0 {
		post.Description = form.Description
	}
	if len(form.FeaturingImagePath) > 0 {
		post.FeaturingImagePath = form.FeaturingImagePath
	}
	if len(form.Categories) > 0 {
		if len(form.realCategories) == 0 {
			return nil, nil, errors.New("validate the form first")
		}
		for _, realCategory := range form.realCategories {
			categories = append(categories, realCategory.ToCommonModel())
		}
		post.Categories = categories
	}
	if len(form.Tags) > 0 {
		post.Tags = form.Tags
	}
	if len(form.Content) > 0 {
		postContent.Content = form.Content
	}
	if form.PublishNow {
		post.PublishedAt = now
	}
	post.UpdatedAt = now

	return post, postContent, nil
}

package forms

import (
	"time"

	"github.com/misterabdul/goblog-server/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreatePostForm struct {
	Slug       string         `json:"slug" binding:"required"`
	Title      string         `json:"title" binding:"required"`
	Categories []postCategory `json:"categories" binding:"required"`
	Tags       []string       `json:"tags"`
	Content    string         `json:"content" binding:"required"`
	PublishNow bool           `json:"publishNow"`
}

type UpdatePostForm struct {
	Slug       string         `json:"slug"`
	Title      string         `json:"title"`
	Categories []postCategory `json:"categories"`
	Tags       []string       `json:"tags"`
	Content    string         `json:"content"`
	PublishNow bool           `json:"publishNow"`
}

type postCategory struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func CreatePostModel(form *CreatePostForm, author *models.UserModel) *models.PostModel {
	var categories []models.CategoryCommonModel
	for _, formCategory := range form.Categories {
		category := models.CategoryCommonModel{
			Name: formCategory.Name,
			Slug: formCategory.Slug,
		}
		categories = append(categories, category)
	}

	now := primitive.NewDateTimeFromTime(time.Now())
	var publishedAt interface{} = nil
	if form.PublishNow {
		publishedAt = now
	}

	return &models.PostModel{
		UID:         primitive.NewObjectID(),
		Slug:        form.Slug,
		Title:       form.Title,
		Categories:  categories,
		Tags:        form.Tags,
		Content:     form.Content,
		Author:      models.CreateUserCommonModel(*author),
		PublishedAt: publishedAt,
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   nil,
	}
}

func UpdatePostModel(form *UpdatePostForm, post *models.PostModel) *models.PostModel {
	now := primitive.NewDateTimeFromTime(time.Now())

	if len(form.Slug) > 0 {
		post.Slug = form.Slug
	}
	if len(form.Title) > 0 {
		post.Title = form.Title
	}
	if len(form.Categories) > 0 {
		var categories []models.CategoryCommonModel
		for _, formCategory := range form.Categories {
			category := models.CategoryCommonModel{
				Name: formCategory.Name,
				Slug: formCategory.Slug,
			}
			categories = append(categories, category)
		}
		post.Categories = categories
	}
	if len(form.Tags) > 0 {
		post.Tags = form.Tags
	}
	if len(form.Content) > 0 {
		post.Content = form.Content
	}
	if form.PublishNow {
		post.PublishedAt = now
	}
	post.UpdatedAt = now

	return post
}
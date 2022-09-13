package forms

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

type CreatePostForm struct {
	Slug               string   `json:"slug" binding:"required,alphanum,max=100"`
	Title              string   `json:"title" binding:"required,max=100"`
	Description        string   `json:"description" binding:"omitempty,max=255"`
	FeaturingImagePath string   `json:"featuringImagePath" binding:"omitempty,url"`
	Categories         []string `json:"categories" binding:"required,dive,len=24"`
	Tags               []string `json:"tags" binding:"omitempty,dive,alphanum,max=32"`
	Content            string   `json:"content" binding:"required"`
	PublishNow         bool     `json:"publishNow" binding:"omitempty"`

	realCategories []*models.CategoryModel
}

func (form *CreatePostForm) Validate(
	categoryService *service.CategoryService,
	postService *service.PostService,
) (err error) {
	if err = checkPostSlug(postService, form.Slug); err != nil {
		return err
	}
	if form.realCategories, err = findCategories(categoryService, form.Categories); err != nil {
		return err
	}
	if len(form.realCategories) == 0 {
		return errors.New("couldn't find any categories from the input")
	}

	return nil
}

func (form *CreatePostForm) ToPostModel(author *models.UserModel) (
	post *models.PostModel,
	content *models.PostContentModel,
	err error,
) {
	var (
		categories              = []models.CategoryCommonModel{}
		now                     = primitive.NewDateTimeFromTime(time.Now())
		postId                  = primitive.NewObjectID()
		publishedAt interface{} = nil
	)

	if len(form.realCategories) == 0 {
		return nil, nil, errors.New("validate the form first")
	}
	for _, realCategory := range form.realCategories {
		categories = append(categories, realCategory.ToCommonModel())
	}
	if form.PublishNow {
		publishedAt = now
	}

	return &models.PostModel{
			UID:                postId,
			Slug:               form.Slug,
			Title:              form.Title,
			Description:        form.Description,
			FeaturingImagePath: form.FeaturingImagePath,
			Categories:         categories,
			Tags:               form.Tags,
			Author:             author.ToCommonModel(),
			PublishedAt:        publishedAt,
			CreatedAt:          now,
			UpdatedAt:          now,
			DeletedAt:          nil,
		}, &models.PostContentModel{
			UID:     postId,
			Content: form.Content}, nil
}

func checkPostSlug(
	postService *service.PostService,
	formSlug string,
) (err error) {
	var (
		posts []*models.PostModel
	)

	if posts, err = postService.GetPosts(bson.M{
		"slug": bson.M{"$eq": formSlug},
	}); err != nil {
		return err
	}
	if len(posts) > 0 {
		return errors.New("slug exists")
	}

	return nil
}

func findCategories(
	categoryService *service.CategoryService,
	formCategories []string,
) (categories []*models.CategoryModel, err error) {
	var categoryUids []primitive.ObjectID

	if categoryUids, err = toObjectIdArray(formCategories); err != nil {
		return nil, err
	}
	if categories, err = categoryService.GetCategories(bson.M{
		"$and": []bson.M{
			{"deletedat": bson.M{"$eq": primitive.Null{}}},
			{"_id": bson.M{"$in": categoryUids}}},
	}); err != nil {
		return nil, err
	}

	return categories, nil
}

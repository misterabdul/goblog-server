package forms

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

type UpdatePageForm struct {
	Slug       string `json:"slug" binding:"omitempty,max=100"`
	Title      string `json:"title" binding:"omitempty,max=100"`
	Content    string `json:"content" binding:"omitempty"`
	PublishNow bool   `json:"publishNow" binding:"omitempty"`
}

func (form *UpdatePageForm) Validate(
	pageService *service.PageService,
	page *models.PageModel,
) (err error) {
	if err = checkUpdatePageSlug(page, pageService, form.Slug); err != nil {
		return err
	}

	return nil
}

func (form *UpdatePageForm) ToPageModel(
	page *models.PageModel,
	pageContent *models.PageContentModel,
) (
	updatedPage *models.PageModel,
	updatedPageContent *models.PageContentModel,
	err error,
) {
	var (
		now = primitive.NewDateTimeFromTime(time.Now())
	)

	if len(form.Slug) > 0 {
		page.Slug = form.Slug
	}
	if len(form.Title) > 0 {
		page.Title = form.Title
	}
	if len(form.Content) > 0 {
		pageContent.Content = form.Content
	}
	if form.PublishNow {
		page.PublishedAt = now
	}
	page.UpdatedAt = now

	return page, pageContent, nil
}

func checkUpdatePageSlug(
	page *models.PageModel,
	pageService *service.PageService,
	formSlug string,
) (err error) {
	var (
		pages []*models.PageModel
	)

	if pages, err = pageService.GetPages(bson.M{
		"$and": []bson.M{
			{"_id": bson.M{"$ne": page.UID}},
			{"slug": bson.M{"$eq": formSlug}}},
	}); err != nil {
		return err
	}
	if len(pages) > 0 {
		return errors.New("slug exists")
	}

	return nil
}

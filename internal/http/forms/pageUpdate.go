package forms

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

type UpdatePageForm struct {
	Slug       string `json:"slug" binding:"omitempty,max=100"`
	Title      string `json:"title" binding:"omitempty,max=100"`
	Content    string `json:"content" binding:"omitempty"`
	PublishNow bool   `json:"publishNow" binding:"omitempty"`
}

func (form *UpdatePageForm) Validate(
	pageService *service.Service,
) (err error) {
	if err = checkPageSlug(pageService, form.Slug); err != nil {
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

package forms

import (
	"errors"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

type CreatePageForm struct {
	Slug       string `json:"slug" binding:"required,max=100"`
	Title      string `json:"title" binding:"required,max=100"`
	Content    string `json:"content" binding:"required"`
	PublishNow bool   `json:"publishNow" binding:"omitempty"`
}

func (form *CreatePageForm) Validate(
	pageService *service.PageService,
) (err error) {
	var (
		parsedUrl *url.URL
	)

	if parsedUrl, err = url.ParseRequestURI(form.Slug); err != nil {
		return err
	}
	form.Slug = parsedUrl.Path
	if err = checkPageSlug(pageService, form.Slug); err != nil {
		return err
	}

	return nil
}

func (form *CreatePageForm) ToPageModel(author *models.UserModel) (
	page *models.PageModel,
	content *models.PageContentModel,
	err error,
) {
	var (
		now                     = primitive.NewDateTimeFromTime(time.Now())
		pageId                  = primitive.NewObjectID()
		publishedAt interface{} = nil
	)

	if form.PublishNow {
		publishedAt = now
	}

	return &models.PageModel{
			UID:         pageId,
			Slug:        form.Slug,
			Title:       form.Title,
			Author:      author.ToCommonModel(),
			PublishedAt: publishedAt,
			CreatedAt:   now,
			UpdatedAt:   now,
			DeletedAt:   nil,
		}, &models.PageContentModel{
			UID:     pageId,
			Content: form.Content}, nil
}

func checkPageSlug(
	pageService *service.PageService,
	formSlug string,
) (err error) {
	var (
		pages []*models.PageModel
	)

	if pages, err = pageService.GetPages(bson.M{
		"slug": bson.M{"$eq": formSlug},
	}); err != nil {
		return err
	}
	if len(pages) > 0 {
		return errors.New("slug exists")
	}

	return nil
}

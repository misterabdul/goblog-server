package pages

import (
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

// @Tags        Page (Public)
// @Summary     Get Public Page
// @Description Get a page that available publicly.
// @Router      /v1/page/{uid} [get]
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Post's UID or slug"
// @Success     200 {object} object{data=object{uid=string,slug=string,title=string,content=string,publishedAt=time}}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func GetPublicPage(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService = service.NewPageService(c, ctx, dbConn)
			page        *models.PageModel
			pageContent *models.PageContentModel
			pageUid     interface{}
			pageParam   = c.Param("page")
			err         error
		)

		defer cancel()
		if pageUid, err = primitive.ObjectIDFromHex(pageParam); err != nil {
			pageUid = nil
		}
		if page, pageContent, err = pageService.GetPageWithContent(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}},
				{"$or": []bson.M{
					{"_id": bson.M{"$eq": pageUid}}}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if page == nil {
			responses.NotFound(c, errors.New("page not found"))
			return
		}

		responses.PublicPage(c, page, pageContent)
	}
}

// @Tags        Page (Public)
// @Summary     Get Public Page By Slug
// @Description Get a page that available publicly by its slug.
// @Router      /v1/page/{slug} [get]
// @Produce     application/json
// @Produce     application/msgpack
// @Param       slug query    string false "The slug query."
// @Success     200  {object} object{data=object{uid=string,slug=string,title=string,content=string,publishedAt=time}}
// @Failure     404  {object} object{message=string}
// @Failure     500  {object} object{message=string}
func GetPublicPageBySlug(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService = service.NewPageService(c, ctx, dbConn)
			page        *models.PageModel
			pageContent *models.PageContentModel
			pageSlug    interface{}
			pageParam   = c.Query("slug")
			err         error
		)

		defer cancel()
		if pageSlug, err = url.Parse(pageParam); err != nil {
			pageSlug = nil
		}
		if page, pageContent, err = pageService.GetPageWithContent(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}},
				{"slug": bson.M{"$eq": pageSlug}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if page == nil {
			responses.NotFound(c, errors.New("page not found"))
			return
		}

		responses.PublicPage(c, page, pageContent)
	}
}

// @Tags        Page (Public)
// @Summary     Get Public Pages
// @Description Get pages that available publicly.
// @Router      /v1/pages [get]
// @Produce     application/json
// @Produce     application/msgpack
// @Param       show  query    int    false "Number of data to be shown."
// @Param       page  query    int    false "Selected page of data."
// @Param       order query    string false "Selected field to order data with."
// @Param       asc   query    string false "Ascending or descending, e.g.: ?asc=false."
// @Success     200   {object} object{data=[]object{uid=string,slug=string,title=string,content=string,publishedAt=time}}
// @Success     204
// @Failure     404   {object} object{message=string}
// @Failure     500   {object} object{message=string}
func GetPublicPages(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService = service.NewPageService(c, ctx, dbConn)
			pages       []*models.PageModel
			err         error
		)

		defer cancel()
		if pages, err = pageService.GetPages(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(pages) == 0 {
			responses.NoContent(c)
			return
		}

		responses.PublicPages(c, pages)
	}
}

// @Tags        Page (Public)
// @Summary     Search Public Pages
// @Description Search posts that available publicly.
// @Router      /v1/page/search [get]
// @Produce     application/json
// @Produce     application/msgpack
// @Param       q     query    string false "The search query."
// @Param       show  query    int    false "Number of data to be shown."
// @Param       page  query    int    false "Selected page of data."
// @Param       order query    string false "Selected field to order data with."
// @Param       asc   query    string false "Ascending or descending, e.g.: ?asc=false."
// @Success     200   {object} object{data=[]object{uid=string,slug=string,title=string,content=string,publishedAt=time}}
// @Success     204
// @Failure     404   {object} object{message=string}
// @Failure     500   {object} object{message=string}
func SearchPublicPages(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService = service.NewPageService(c, ctx, dbConn)
			searchQuery = c.Query("q")
			pages       []*models.PageModel
			err         error
		)

		defer cancel()
		if pages, err = pageService.GetPages(bson.M{
			"$text": bson.M{"$search": searchQuery},
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(pages) == 0 {
			responses.NoContent(c)
			return
		}

		responses.PublicPages(c, pages)
	}
}

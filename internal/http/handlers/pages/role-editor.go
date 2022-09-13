package pages

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/service"
)

// @Tags        Page (Editor)
// @Summary     Get Page
// @Description Get a page.
// @Router      /v1/auth/editor/page/{uid} [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Page's UID or slug"
// @Success     200 {object} object{data=object{uid=string,slug=string,title=string,content=string,author=object{uid=string,username=string,email=string,firstName=string,lastName=string},publishedAt=time,updatedAt=time,createdAt=time,deletedAt=time}}
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func GetPage(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService  = service.NewPageService(c, ctx, dbConn)
			page         *models.PageModel
			pageContent  *models.PageContentModel
			pageUid      primitive.ObjectID
			pageUidParam = c.Param("page")
			err          error
		)

		defer cancel()
		if pageUid, err = primitive.ObjectIDFromHex(pageUidParam); err != nil {
			responses.IncorrectPageId(c, err)
			return
		}
		if page, pageContent, err = pageService.GetPageWithContent(bson.M{
			"_id": bson.M{"$eq": pageUid},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if page == nil {
			responses.NotFound(c, errors.New("page not found"))
			return
		}

		responses.AuthorizedPage(c, page, pageContent)
	}
}

// @Tags        Page (Editor)
// @Summary     Get Pages
// @Description Get pages.
// @Router      /v1/auth/editor/pages [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       show  query    int    false "Number of data to be shown."
// @Param       page  query    int    false "Selected page of data."
// @Param       order query    string false "Selected field to order data with."
// @Param       asc   query    string false "Ascending or descending, e.g.: ?asc=false."
// @Param       type  query    string false "Filter data by type, e.g.: ?type=trash, ?type=published, ?type=draft."
// @Success     200   {object} object{data=[]object{uid=string,slug=string,title=string,content=string,author=object{uid=string,username=string,email=string,firstName=string,lastName=string},publishedAt=time,updatedAt=time,createdAt=time,deletedAt=time}}
// @Success     204
// @Failure     401   {object} object{message=string}
// @Failure     500   {object} object{message=string}
func GetPages(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService = service.NewPageService(c, ctx, dbConn)
			pages       []*models.PageModel
			queryParams = readCommonQueryParams(c)
			err         error
		)

		defer cancel()
		if pages, err = pageService.GetPages(bson.M{
			"$and": queryParams,
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(pages) == 0 {
			responses.NoContent(c)
			return
		}

		responses.AuthorizedPages(c, pages)
	}
}

// @Tags        Page (Editor)
// @Summary     Get Pages Stats
// @Description Get pages's stats.
// @Router      /v1/auth/editor/pages/stats [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       show  query    int    false "Number of data to be shown."
// @Param       page  query    int    false "Selected page of data."
// @Param       order query    string false "Selected field to order data with."
// @Param       asc   query    string false "Ascending or descending, e.g.: ?asc=false."
// @Param       type  query    string false "Filter data by type, e.g.: ?type=trash, ?type=published, ?type=draft."
// @Success     200   {object} object{data=object{currentPage=int,totalPages=int,itemsPerPage=int,totalItems=int}}
// @Failure     401   {object} object{message=string}
// @Failure     500   {object} object{message=string}
func GetPagesStats(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService = service.NewPageService(c, ctx, dbConn)
			count       int64
			queryParams = readCommonQueryParams(c)
			err         error
		)

		defer cancel()
		if count, err = pageService.GetPageCount(bson.M{
			"$and": queryParams,
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.ResourceStats(c, count)
	}
}

// @Tags        Page (Editor)
// @Summary     Create Page
// @Description Create a new page.
// @Router      /v1/auth/editor/page [post]
// @Security    BearerAuth
// @Accept      application/json
// @Accept      application/msgpack
// @Produce     application/json
// @Produce     application/msgpack
// @Param       form body     object{slug=string,title=string,content=string,publishNow=boolean} true "Create page form"
// @Success     200  {object} object{data=object{uid=string,slug=string,title=string,content=string,author=object{uid=string,username=string,email=string,firstName=string,lastName=string},publishedAt=time,updatedAt=time,createdAt=time,deletedAt=time}}
// @Failure     401  {object} object{message=string}
// @Failure     422  {object} object{message=string}
// @Failure     500  {object} object{message=string}
func CreatePage(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService = service.NewPageService(c, ctx, dbConn)
			me          *models.UserModel
			page        *models.PageModel
			pageContent *models.PageContentModel
			form        *forms.CreatePageForm
			err         error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if form, err = requests.GetCreatePageForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if err = form.Validate(pageService); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if page, pageContent, err = form.ToPageModel(me); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = pageService.CreatePage(page, pageContent); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.MyPage(c, page, pageContent)
	}
}

// @Tags        Page (Editor)
// @Summary     Publish Page
// @Description Publish a page if not published yet.
// @Router      /v1/auth/editor/page/{uid}/publish [put]
// @Router      /v1/auth/editor/page/{uid}/publish [patch]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func PublishPage(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService  = service.NewPageService(c, ctx, dbConn)
			page         *models.PageModel
			pageUid      primitive.ObjectID
			pageUidParam = c.Param("page")
			err          error
		)

		defer cancel()
		if pageUid, err = primitive.ObjectIDFromHex(pageUidParam); err != nil {
			responses.IncorrectPageId(c, err)
			return
		}
		if page, err = pageService.GetPage(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"_id": bson.M{"$eq": pageUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if page == nil {
			responses.NotFound(c, err)
			return
		}
		if page.PublishedAt != nil {
			responses.NoContent(c)
			return
		}
		if err = pageService.PublishPage(page); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        Page (Editor)
// @Summary     Unpublish Page
// @Description Remove publish status from a page if page already published.
// @Router      /v1/auth/writer/page/{uid}/depublish [put]
// @Router      /v1/auth/writer/page/{uid}/depublish [patch]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func DepublishPage(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService  = service.NewPageService(c, ctx, dbConn)
			page         *models.PageModel
			pageUid      primitive.ObjectID
			pageUidParam = c.Param("page")
			err          error
		)

		defer cancel()
		if pageUid, err = primitive.ObjectIDFromHex(pageUidParam); err != nil {
			responses.IncorrectPageId(c, err)
			return
		}
		if page, err = pageService.GetPage(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"_id": bson.M{"$eq": pageUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if page == nil {
			responses.NotFound(c, errors.New("page not found"))
			return
		}
		if page.PublishedAt == nil {
			responses.NoContent(c)
			return
		}
		if err = pageService.DepublishPage(page); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        Page (Editor)
// @Summary     Update Page
// @Description Update a page.
// @Router      /v1/auth/editor/page/{uid} [put]
// @Router      /v1/auth/editor/page/{uid} [patch]
// @Security    BearerAuth
// @Accept      application/json
// @Accept      application/msgpack
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid  path     string                                                             true "Page's UID or slug"
// @Param       form body     object{slug=string,title=string,content=string,publishNow=boolean} true "Update page form"
// @Success     204
// @Failure     401  {object} object{message=string}
// @Failure     404  {object} object{message=string}
// @Failure     422  {object} object{message=string}
// @Failure     500  {object} object{message=string}
func UpdatePage(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel        = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService        = service.NewPageService(c, ctx, dbConn)
			page               *models.PageModel
			updatedPage        *models.PageModel
			pageContent        *models.PageContentModel
			updatedPageContent *models.PageContentModel
			pageUid            primitive.ObjectID
			pageUidParam       = c.Param("page")
			form               *forms.UpdatePageForm
			err                error
		)

		defer cancel()
		if pageUid, err = primitive.ObjectIDFromHex(pageUidParam); err != nil {
			responses.IncorrectPageId(c, err)
			return
		}
		if page, pageContent, err = pageService.GetPageWithContent(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"_id": bson.M{"$eq": pageUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if page == nil {
			responses.NotFound(c, errors.New("page not found"))
			return
		}
		if form, err = requests.GetUpdatePageForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if err = form.Validate(pageService, page); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if updatedPage, updatedPageContent, err = form.ToPageModel(
			page, pageContent,
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = pageService.UpdatePage(
			updatedPage, updatedPageContent,
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        Page (Editor)
// @Summary     Delete Page (Soft)
// @Description Delete a page (soft-deleted).
// @Router      /v1/auth/editor/page/{uid} [delete]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Page's UID or slug"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func TrashPage(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService  = service.NewPageService(c, ctx, dbConn)
			page         *models.PageModel
			pageUid      primitive.ObjectID
			pageUidParam = c.Param("page")
			err          error
		)

		defer cancel()
		if pageUid, err = primitive.ObjectIDFromHex(pageUidParam); err != nil {
			responses.IncorrectPageId(c, err)
			return
		}
		if page, err = pageService.GetPage(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"_id": bson.M{"$eq": pageUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if page == nil {
			responses.NotFound(c, errors.New("page not found"))
			return
		}
		if page.DeletedAt != nil {
			responses.NoContent(c)
			return
		}
		if err = pageService.TrashPage(page); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        Page (Editor)
// @Summary     Restore Page (Soft)
// @Description Restore a deleted page (soft-deleted).
// @Router      /v1/auth/editor/page/{uid}/detrash [put]
// @Router      /v1/auth/editor/page/{uid}/detrash [patch]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Page's UID or slug"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func DetrashPage(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService  = service.NewPageService(c, ctx, dbConn)
			page         *models.PageModel
			pageUid      primitive.ObjectID
			pageUidParam = c.Param("page")
			err          error
		)

		defer cancel()
		if pageUid, err = primitive.ObjectIDFromHex(pageUidParam); err != nil {
			responses.IncorrectPageId(c, err)
			return
		}
		if page, err = pageService.GetPage(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": bson.M{"$eq": pageUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if page == nil {
			responses.NotFound(c, errors.New("page not found"))
			return
		}
		if err = pageService.DetrashPage(page); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        Page (Editor)
// @Summary     Delete Page (Permanent)
// @Description Delete a page (permanent).
// @Router      /v1/auth/editor/page/{uid}/permanent [delete]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Page's UID or slug"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func DeletePage(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService  = service.NewPageService(c, ctx, dbConn)
			page         *models.PageModel
			pageContent  *models.PageContentModel
			pageUid      primitive.ObjectID
			pageUidParam = c.Param("page")
			err          error
		)

		defer cancel()
		if pageUid, err = primitive.ObjectIDFromHex(pageUidParam); err != nil {
			responses.IncorrectPageId(c, err)
			return
		}
		if page, pageContent, err = pageService.GetPageWithContent(bson.M{
			"_id": bson.M{"$eq": pageUid},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if page == nil {
			responses.NotFound(c, errors.New("page not found"))
			return
		}
		if err = pageService.DeletePage(page, pageContent); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func readCommonQueryParams(c *gin.Context) []bson.M {
	var (
		typeQuery  = c.DefaultQuery("type", "draft")
		extraQuery = []bson.M{}
	)

	switch true {
	case typeQuery == "trash":
		extraQuery = append(extraQuery,
			bson.M{"deletedat": bson.M{"$ne": primitive.Null{}}})
	case typeQuery == "published":
		extraQuery = append(extraQuery,
			bson.M{"publishedat": bson.M{"$ne": primitive.Null{}}},
			bson.M{"deletedat": bson.M{"$eq": primitive.Null{}}})
	case typeQuery == "draft":
		fallthrough
	default:
		extraQuery = append(extraQuery,
			bson.M{"publishedat": bson.M{"$eq": primitive.Null{}}},
			bson.M{"deletedat": bson.M{"$eq": primitive.Null{}}})
	}

	return extraQuery
}

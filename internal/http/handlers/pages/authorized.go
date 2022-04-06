package pages

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

func GetPage(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService  = service.New(c, ctx, dbConn)
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

func GetPages(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService = service.New(c, ctx, dbConn)
			pages       []*models.PageModel
			typeQuery   = c.DefaultQuery("type", "draft")
			extraQuery  = []bson.M{}
			err         error
		)

		defer cancel()
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
		if pages, err = pageService.GetPages(bson.M{
			"$and": extraQuery,
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

func CreatePage(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService = service.New(c, ctx, dbConn)
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

func PublishPage(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService  = service.New(c, ctx, dbConn)
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

func DepublishPage(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService  = service.New(c, ctx, dbConn)
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

func UpdatePage(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel        = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService        = service.New(c, ctx, dbConn)
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

func TrashPage(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService  = service.New(c, ctx, dbConn)
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

func DetrashPage(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService  = service.New(c, ctx, dbConn)
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

func DeletePage(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService  = service.New(c, ctx, dbConn)
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

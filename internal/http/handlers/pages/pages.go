package pages

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

func GetPublicPage(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService = service.New(c, ctx, dbConn)
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
					{"_id": bson.M{"$eq": pageUid}},
					{"slug": bson.M{"$eq": pageParam}}}}},
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

func GetPublicPages(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService = service.New(c, ctx, dbConn)
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

func SearchPublicPages(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			pageService = service.New(c, ctx, dbConn)
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

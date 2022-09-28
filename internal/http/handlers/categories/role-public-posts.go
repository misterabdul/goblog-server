package categories

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	"github.com/misterabdul/goblog-server/internal/service"
)

// @Tags        Category (Public)
// @Summary     Get Public Category Posts
// @Description Get public category's posts that available publicly.
// @Router      /v1/category/{uid}/posts [get]
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Category's UID or slug"
// @Success     200 {object} object{data=object{uid=string,slug=string,title=string,featuringImagePath=string,description=string,categories=[]object{uid=string,slug=string,name=string},tags=[]string,content=string,author=object{uid=string,username=string,email=string,firstName=string,lastName=string},commentCount=int,publishedAt=time}}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func GetPublicCategoryPosts(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel   = context.WithTimeout(context.Background(), maxCtxDuration)
			posts         []*models.PostModel
			categoryUid   interface{}
			categoryParam = c.Param("category")
			err           error
		)

		defer cancel()
		if categoryUid, err = primitive.ObjectIDFromHex(categoryParam); err != nil {
			categoryUid = nil
		}
		if posts, err = svc.Post.GetMany(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}},
				{"$or": []bson.M{
					{"_id": bson.M{"$eq": categoryUid}},
					{"slug": bson.M{"$eq": categoryParam}}}}}},
			internalGin.GetFindOptions(c),
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(posts) == 0 {
			responses.NoContent(c)
			return
		}

		responses.PublicPosts(c, posts)
	}
}

package posts

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	"github.com/misterabdul/goblog-server/internal/service"
)

// @Tags        Post (Public)
// @Summary     Get Public Post
// @Description Get a post that available publicly.
// @Router      /v1/post/{uid} [get]
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Post's UID or slug"
// @Success     200 {object} object{data=object{uid=string,slug=string,title=string,featuringImagePath=string,description=string,categories=[]object{uid=string,slug=string,name=string},tags=[]string,content=string,author=object{uid=string,username=string,email=string,firstName=string,lastName=string},commentCount=int,publishedAt=time}}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func GetPublicPost(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			post        *models.PostModel
			postContent *models.PostContentModel
			postUid     interface{}
			postParam   = c.Param("post")
			err         error
		)

		defer cancel()
		if postUid, err = primitive.ObjectIDFromHex(postParam); err != nil {
			postUid = nil
		}
		if post, postContent, err = svc.Post.GetOneWithContent(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}},
				{"$or": []bson.M{
					{"_id": bson.M{"$eq": postUid}},
					{"slug": bson.M{"$eq": postParam}}}}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}

		responses.PublicPost(c, post, postContent)
	}
}

// @Tags        Post (Public)
// @Summary     Get Public Posts
// @Description Get posts that available publicly.
// @Router      /v1/posts [get]
// @Produce     application/json
// @Produce     application/msgpack
// @Param       show  query    int    false "Number of data to be shown."
// @Param       page  query    int    false "Selected page of data."
// @Param       order query    string false "Selected field to order data with."
// @Param       asc   query    string false "Ascending or descending, e.g.: ?asc=false."
// @Success     200   {object} object{data=[]object{uid=string,slug=string,title=string,featuringImagePath=string,description=string,categories=[]object{uid=string,slug=string,name=string},tags=[]string,content=string,author=object{uid=string,username=string,email=string,firstName=string,lastName=string},commentCount=int,publishedAt=time}}
// @Success     204
// @Failure     404   {object} object{message=string}
// @Failure     500   {object} object{message=string}
func GetPublicPosts(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			posts       []*models.PostModel
			err         error
		)

		defer cancel()
		if posts, err = svc.Post.GetMany(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}}}},
			internalGin.GetFindOptionsPost(c),
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

// @Tags        Post (Public)
// @Summary     Search Public Posts
// @Description Search posts that available publicly.
// @Router      /v1/post/search [get]
// @Produce     application/json
// @Produce     application/msgpack
// @Param       q     query    string false "The search query."
// @Param       show  query    int    false "Number of data to be shown."
// @Param       page  query    int    false "Selected page of data."
// @Param       order query    string false "Selected field to order data with."
// @Param       asc   query    string false "Ascending or descending, e.g.: ?asc=false."
// @Success     200   {object} object{data=[]object{uid=string,slug=string,title=string,featuringImagePath=string,description=string,categories=[]object{uid=string,slug=string,name=string},tags=[]string,content=string,author=object{uid=string,username=string,email=string,firstName=string,lastName=string},commentCount=int,publishedAt=time}}
// @Success     204
// @Failure     404   {object} object{message=string}
// @Failure     500   {object} object{message=string}
func SearchPublicPosts(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			searchQuery = c.Query("q")
			posts       []*models.PostModel
			err         error
		)

		defer cancel()
		if posts, err = svc.Post.GetMany(ctx, bson.M{
			"$text": bson.M{"$search": searchQuery},
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}}}},
			internalGin.GetFindOptionsPost(c),
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

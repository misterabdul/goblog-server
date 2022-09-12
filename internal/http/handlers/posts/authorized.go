package posts

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

// @Tags        Post (Editor)
// @Summary     Get Post
// @Description Get a post.
// @Router      /v1/auth/editor/post/{uid} [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Post's UID or slug"
// @Success     200 {object} object{data=object{uid=string,slug=string,title=string,featuringImagePath=string,description=string,categories=[]object{uid=string,slug=string,name=string},tags=[]string,author=object{uid=string,username=string,email=string,firstName=string,lastName=string},commentCount=int,publishedAt=time,updatedAt=time,createdAt=time,deletedAt=time}}
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func GetPost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			postService  = service.NewPostService(c, ctx, dbConn)
			post         *models.PostModel
			postContent  *models.PostContentModel
			postUid      primitive.ObjectID
			postUidParam = c.Param("post")
			err          error
		)

		defer cancel()
		if postUid, err = primitive.ObjectIDFromHex(postUidParam); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, postContent, err = postService.GetPostWithContent(bson.M{
			"_id": bson.M{"$eq": postUid},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}

		responses.AuthorizedPost(c, post, postContent)
	}
}

// @Tags        Post (Editor)
// @Summary     Get Posts
// @Description Get posts.
// @Router      /v1/auth/editor/posts [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       show  query    int    false "Number of data to be shown."
// @Param       page  query    int    false "Selected page of data."
// @Param       order query    string false "Selected field to order data with."
// @Param       asc   query    string false "Ascending or descending, e.g.: ?asc=false."
// @Param       type  query    string false "Filter data by type, e.g.: ?type=trash, ?type=published, ?type=draft."
// @Success     200   {object} object{data=[]object{uid=string,slug=string,title=string,featuringImagePath=string,description=string,categories=[]object{uid=string,slug=string,name=string},tags=[]string,author=object{uid=string,username=string,email=string,firstName=string,lastName=string},commentCount=int,publishedAt=time,updatedAt=time,createdAt=time,deletedAt=time}}
// @Success     204
// @Failure     401   {object} object{message=string}
// @Failure     500   {object} object{message=string}
func GetPosts(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			postService = service.NewPostService(c, ctx, dbConn)
			posts       []*models.PostModel
			queryParams = readCommonQueryParams(c)
			err         error
		)

		defer cancel()
		if posts, err = postService.GetPosts(bson.M{
			"$and": queryParams,
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(posts) == 0 {
			responses.NoContent(c)
			return
		}

		responses.AuthorizedPosts(c, posts)
	}
}

// @Tags        Post (Editor)
// @Summary     Get Posts Stats
// @Description Get posts's stats.
// @Router      /v1/auth/editor/posts/stats [get]
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
func GetPostsStats(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			postService = service.NewPostService(c, ctx, dbConn)
			count       int64
			queryParams = readCommonQueryParams(c)
			err         error
		)

		defer cancel()
		if count, err = postService.GetPostCount(bson.M{
			"$and": queryParams,
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.ResourceStats(c, count)
	}
}

// @Tags        Post (Editor)
// @Summary     Publish Post
// @Description Publish a post if not published yet.
// @Router      /v1/auth/editor/post/{uid}/publish [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func PublishPost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			postService  = service.NewPostService(c, ctx, dbConn)
			post         *models.PostModel
			postUid      primitive.ObjectID
			postUidParam = c.Param("post")
			err          error
		)

		defer cancel()
		if postUid, err = primitive.ObjectIDFromHex(postUidParam); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, err = postService.GetPost(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"_id": bson.M{"$eq": postUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, err)
			return
		}
		if post.PublishedAt != nil {
			responses.NoContent(c)
			return
		}
		if err = postService.PublishPost(post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        Post (Editor)
// @Summary     Unpublish Post
// @Description Remove a post from publish status if post already published.
// @Router      /v1/auth/editor/post/{uid}/depublish [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func DepublishPost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			postService  = service.NewPostService(c, ctx, dbConn)
			post         *models.PostModel
			postUid      primitive.ObjectID
			postUidParam = c.Param("post")
			err          error
		)

		defer cancel()
		if postUid, err = primitive.ObjectIDFromHex(postUidParam); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, err = postService.GetPost(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"_id": bson.M{"$eq": postUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if post.PublishedAt == nil {
			responses.NoContent(c)
			return
		}
		if err = postService.DepublishPost(post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        Post (Editor)
// @Summary     Update Post
// @Description Update a post.
// @Router      /v1/auth/editor/post/{uid} [put]
// @Router      /v1/auth/editor/post/{uid} [patch]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid  path string                                                                                                                                            true "Post's UID or slug"
// @Param       form body object{slug=string,title=string,description=string,featuringImagePath=string,categories=[]string,tags=[]string,content=string,publishNow=boolean} true "Update post form"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     422 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func UpdatePost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel        = context.WithTimeout(context.Background(), maxCtxDuration)
			postService        = service.NewPostService(c, ctx, dbConn)
			categoryService    = service.NewCategoryService(c, ctx, dbConn)
			post               *models.PostModel
			updatedPost        *models.PostModel
			postContent        *models.PostContentModel
			updatedPostContent *models.PostContentModel
			postUid            primitive.ObjectID
			postUidParam       = c.Param("post")
			form               *forms.UpdatePostForm
			err                error
		)

		defer cancel()
		if postUid, err = primitive.ObjectIDFromHex(postUidParam); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, postContent, err = postService.GetPostWithContent(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"_id": bson.M{"$eq": postUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if form, err = requests.GetUpdatePostForm(c); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if err = form.Validate(categoryService, postService, post); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if updatedPost, updatedPostContent, err = form.ToPostModel(
			post, postContent,
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = postService.UpdatePost(
			updatedPost, updatedPostContent,
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        Post (Editor)
// @Summary     Delete Post (Soft)
// @Description Delete a post (soft-deleted).
// @Router      /v1/auth/editor/post/{uid} [delete]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path string true "Post's UID or slug"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     422 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func TrashPost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			postService  = service.NewPostService(c, ctx, dbConn)
			post         *models.PostModel
			postUid      primitive.ObjectID
			postUidParam = c.Param("post")
			err          error
		)

		defer cancel()
		if postUid, err = primitive.ObjectIDFromHex(postUidParam); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, err = postService.GetPost(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"_id": bson.M{"$eq": postUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if post.DeletedAt != nil {
			responses.NoContent(c)
			return
		}
		if err = postService.TrashPost(post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        Post (Editor)
// @Summary     Restore Post (Soft)
// @Description Restore a deleted post (soft-deleted).
// @Router      /v1/auth/editor/post/{uid}/detrash [put]
// @Router      /v1/auth/editor/post/{uid}/detrash [patch]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path string true "Post's UID or slug"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     422 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func DetrashPost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			postService  = service.NewPostService(c, ctx, dbConn)
			post         *models.PostModel
			postUid      primitive.ObjectID
			postUidParam = c.Param("post")
			err          error
		)

		defer cancel()
		if postUid, err = primitive.ObjectIDFromHex(postUidParam); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, err = postService.GetPost(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": bson.M{"$eq": postUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if err = postService.DetrashPost(post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        Post (Writer)
// @Summary     Delete Post (Permanent)
// @Description Delete a post (permanent).
// @Router      /v1/auth/editor/post/{uid}/permanent [delete]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Post's UID or slug"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     422 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func DeletePost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			postService  = service.NewPostService(c, ctx, dbConn)
			post         *models.PostModel
			postContent  *models.PostContentModel
			postUid      primitive.ObjectID
			postUidParam = c.Param("post")
			err          error
		)

		defer cancel()
		if postUid, err = primitive.ObjectIDFromHex(postUidParam); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, postContent, err = postService.GetPostWithContent(bson.M{
			"_id": bson.M{"$eq": postUid},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if err = postService.DeletePost(post, postContent); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

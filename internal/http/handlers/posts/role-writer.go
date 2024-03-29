package posts

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	"github.com/misterabdul/goblog-server/internal/service"
)

// @Tags        Post (Writer)
// @Summary     Get My Post
// @Description Get my post.
// @Router      /v1/auth/writer/post/{uid} [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Post's UID or slug"
// @Success     200 {object} object{data=object{uid=string,slug=string,title=string,featuringImagePath=string,description=string,categories=[]object{uid=string,slug=string,name=string},tags=[]string,author=object{uid=string,username=string,email=string,firstName=string,lastName=string},commentCount=int,publishedAt=time,updatedAt=time,createdAt=time,deletedAt=time}}
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func GetMyPost(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			me           *models.UserModel
			post         *models.PostModel
			postContent  *models.PostContentModel
			postUid      primitive.ObjectID
			postUidParam = c.Param("post")
			err          error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if postUid, err = primitive.ObjectIDFromHex(postUidParam); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, postContent, err = svc.Post.GetOneWithContent(ctx, bson.M{
			"$and": []bson.M{
				{"author._id": bson.M{"$eq": me.UID}},
				{"_id": bson.M{"$eq": postUid}}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}

		responses.MyPost(c, post, postContent)
	}
}

// @Tags        Post (Writer)
// @Summary     Get My Posts
// @Description Get my posts.
// @Router      /v1/auth/writer/posts [get]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       show  query    int    false "Number of data to be shown."
// @Param       page  query    int    false "Selected page of data."
// @Param       order query    string false "Selected field to order data with."
// @Param       asc   query    string false "Ascending or descending, e.g.: ?asc=false."
// @Param       type  query    string false "Filter data by type, e.g.: ?type=trash, ?type=published, ?type=draft."
// @Success     200   {object} object{data=[]object{uid=string,slug=string,title=string,featuringImagePath=string,description=string,categories=[]object{uid=string,slug=string,name=string},tags=[]string,author=object{uid=string,username=string,email=string,firstName=string,lastName=string},commentCount=int,publishedAt=time,updatedAt=time,createdAt=time,deletedAt=time}}
// @Failure     204
// @Failure     401   {object} object{message=string}
// @Failure     500   {object} object{message=string}
func GetMyPosts(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			me          *models.UserModel
			posts       []*models.PostModel
			queryParams = readCommonQueryParams(c)
			err         error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if posts, err = svc.Post.GetMany(ctx, bson.M{
			"$and": append(queryParams,
				bson.M{"author._id": bson.M{"$eq": me.UID}})},
			internalGin.GetFindOptionsPost(c),
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(posts) == 0 {
			responses.NoContent(c)
			return
		}

		responses.MyPosts(c, posts)
	}
}

// @Tags        Post (Writer)
// @Summary     Get My Posts Stats
// @Description Get my posts's stats.
// @Router      /v1/auth/writer/posts/stats [get]
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
func GetMyPostsStats(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			me          *models.UserModel
			count       int64
			queryParams = readCommonQueryParams(c)
			err         error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if count, err = svc.Post.Count(ctx, bson.M{
			"$and": append(queryParams,
				bson.M{"author._id": bson.M{"$eq": me.UID}})},
			internalGin.GetCountOptions(c),
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.ResourceStats(c, count)
	}
}

// @Tags        Post (Writer)
// @Summary     Create Post
// @Description Create a new post.
// @Router      /v1/auth/writer/post [post]
// @Security    BearerAuth
// @Accept      application/json
// @Accept      application/msgpack
// @Produce     application/json
// @Produce     application/msgpack
// @Param       form body     object{slug=string,title=string,description=string,featuringImagePath=string,categories=[]string,tags=[]string,content=string,publishNow=boolean} true "Create post form"
// @Success     200  {object} object{data=object{uid=string,slug=string,title=string,featuringImagePath=string,description=string,categories=[]object{uid=string,slug=string,name=string},tags=[]string,author=object{uid=string,username=string,email=string,firstName=string,lastName=string},commentCount=int,publishedAt=time,updatedAt=time,createdAt=time,deletedAt=time}}
// @Failure     401  {object} object{message=string}
// @Failure     422  {object} object{message=string}
// @Failure     500  {object} object{message=string}
func CreatePost(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			me          *models.UserModel
			post        *models.PostModel
			postContent *models.PostContentModel
			form        *forms.CreatePostForm
			err         error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if form, err = requests.GetCreatePostForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if err = form.Validate(svc, ctx); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if post, postContent, err = form.ToPostModel(me); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = svc.Post.SaveOneWithContent(ctx, post, postContent); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.MyPost(c, post, postContent)
	}
}

// @Tags        Post (Writer)
// @Summary     Publish My Post
// @Description Publish my post if not published yet.
// @Router      /v1/auth/writer/post/{uid}/publish [put]
// @Router      /v1/auth/writer/post/{uid}/publish [patch]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func PublishMyPost(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			me           *models.UserModel
			post         *models.PostModel
			postUid      primitive.ObjectID
			postUidParam = c.Param("post")
			err          error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if postUid, err = primitive.ObjectIDFromHex(postUidParam); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, err = svc.Post.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"author._id": bson.M{"$eq": me.UID}},
				{"_id": bson.M{"$eq": postUid}}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if post.PublishedAt != nil {
			responses.NoContent(c)
			return
		}
		if err = svc.Post.PublishOne(ctx, post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        Post (Writer)
// @Summary     Unpublish My Post
// @Description Remove publish status from my post if post already published.
// @Router      /v1/auth/writer/post/{uid}/depublish [put]
// @Router      /v1/auth/writer/post/{uid}/depublish [patch]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func DepublishMyPost(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			me           *models.UserModel
			post         *models.PostModel
			postUid      primitive.ObjectID
			postUidParam = c.Param("post")
			err          error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if postUid, err = primitive.ObjectIDFromHex(postUidParam); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, err = svc.Post.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"author._id": bson.M{"$eq": me.UID}},
				{"_id": bson.M{"$eq": postUid}}}},
		); err != nil {
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
		if err = svc.Post.DepublishOne(ctx, post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        Post (Writer)
// @Summary     Update My Post
// @Description Update my post.
// @Router      /v1/auth/writer/post/{uid} [put]
// @Router      /v1/auth/writer/post/{uid} [patch]
// @Security    BearerAuth
// @Accept      application/json
// @Accept      application/msgpack
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid  path     string                                                                                                                                            true "Post's UID or slug"
// @Param       form body     object{slug=string,title=string,description=string,featuringImagePath=string,categories=[]string,tags=[]string,content=string,publishNow=boolean} true "Update post form"
// @Success     204
// @Failure     401  {object} object{message=string}
// @Failure     404  {object} object{message=string}
// @Failure     422  {object} object{message=string}
// @Failure     500  {object} object{message=string}
func UpdateMyPost(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel        = context.WithTimeout(context.Background(), maxCtxDuration)
			me                 *models.UserModel
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
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if postUid, err = primitive.ObjectIDFromHex(postUidParam); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, postContent, err = svc.Post.GetOneWithContent(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"author._id": bson.M{"$eq": me.UID}},
				{"_id": bson.M{"$eq": postUid}}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if form, err = requests.GetUpdatePostForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if err = form.Validate(svc, ctx, post); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if updatedPost, updatedPostContent, err = form.ToPostModel(post, postContent); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if err = svc.Post.UpdateOneWithContent(ctx, updatedPost, updatedPostContent); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        Post (Writer)
// @Summary     Delete My Post (Soft)
// @Description Delete my post (soft-deleted).
// @Router      /v1/auth/writer/post/{uid} [delete]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Post's UID or slug"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func TrashMyPost(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			me           *models.UserModel
			post         *models.PostModel
			postUid      primitive.ObjectID
			postUidParam = c.Param("post")
			err          error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if postUid, err = primitive.ObjectIDFromHex(postUidParam); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, err = svc.Post.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"author._id": bson.M{"$eq": me.UID}},
				{"_id": bson.M{"$eq": postUid}}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if err = svc.Post.TrashOne(ctx, post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        Post (Writer)
// @Summary     Restore My Post (Soft)
// @Description Restore my deleted post (soft-deleted).
// @Router      /v1/auth/writer/post/{uid}/detrash [put]
// @Router      /v1/auth/writer/post/{uid}/detrash [patch]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Post's UID or slug"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func DetrashMyPost(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			me           *models.UserModel
			post         *models.PostModel
			postUid      primitive.ObjectID
			postUidParam = c.Param("post")
			err          error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if postUid, err = primitive.ObjectIDFromHex(postUidParam); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, err = svc.Post.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"author._id": bson.M{"$eq": me.UID}},
				{"_id": bson.M{"$eq": postUid}}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, err)
			return
		}
		if err = svc.Post.RestoreOne(ctx, post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

// @Tags        Post (Writer)
// @Summary     Delete My Post (Permanent)
// @Description Delete my post (permanent).
// @Router      /v1/auth/writer/post/{uid}/permanent [delete]
// @Security    BearerAuth
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Post's UID or slug"
// @Success     204
// @Failure     401 {object} object{message=string}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func DeleteMyPost(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			me           *models.UserModel
			post         *models.PostModel
			postContent  *models.PostContentModel
			postUid      primitive.ObjectID
			postUidParam = c.Param("post")
			err          error
		)

		defer cancel()
		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if postUid, err = primitive.ObjectIDFromHex(postUidParam); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, postContent, err = svc.Post.GetOneWithContent(ctx, bson.M{
			"$and": []bson.M{
				{"author._id": bson.M{"$eq": me.UID}},
				{"_id": bson.M{"$eq": postUid}}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, err)
			return
		}
		if err = svc.Post.DeleteOneWithContent(ctx, post, postContent); err != nil {
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

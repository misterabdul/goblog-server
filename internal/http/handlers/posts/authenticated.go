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
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/service"
)

func GetMyPost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			postService  = service.New(c, ctx, dbConn)
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
		if post, postContent, err = postService.GetPostWithContent(bson.M{
			"$and": []bson.M{
				{"author._id": bson.M{"$eq": me.UID}},
				{"_id": bson.M{"$eq": postUid}}},
		}); err != nil {
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

func GetMyPosts(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			postService = service.New(c, ctx, dbConn)
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
		if posts, err = postService.GetPosts(bson.M{
			"$and": append(queryParams,
				bson.M{"author._id": bson.M{"$eq": me.UID}}),
		}); err != nil {
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

func GetMyPostsStats(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			postService = service.New(c, ctx, dbConn)
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
		if count, err = postService.GetPostCount(bson.M{
			"$and": append(queryParams,
				bson.M{"author._id": bson.M{"$eq": me.UID}}),
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.ResourceStats(c, count)
	}
}

func CreatePost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			postService = service.New(c, ctx, dbConn)
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
		if err = form.Validate(postService); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if post, postContent, err = form.ToPostModel(me); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = postService.CreatePost(post, postContent); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.MyPost(c, post, postContent)
	}
}

func PublishMyPost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			postService  = service.New(c, ctx, dbConn)
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
		if post, err = postService.GetPost(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"author._id": bson.M{"$eq": me.UID}},
				{"_id": bson.M{"$eq": postUid}}},
		}); err != nil {
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
		if err = postService.PublishPost(post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func DepublishMyPost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			postService  = service.New(c, ctx, dbConn)
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
		if post, err = postService.GetPost(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"author._id": bson.M{"$eq": me.UID}},
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

func UpdateMyPost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel        = context.WithTimeout(context.Background(), maxCtxDuration)
			postService        = service.New(c, ctx, dbConn)
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
		if post, postContent, err = postService.GetPostWithContent(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"author._id": bson.M{"$eq": me.UID}},
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
			responses.FormIncorrect(c, err)
			return
		}
		if err = form.Validate(postService, post); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if updatedPost, updatedPostContent, err = form.ToPostModel(
			post, postContent,
		); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if err = postService.UpdatePost(
			updatedPost,
			updatedPostContent,
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func TrashMyPost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			postService  = service.New(c, ctx, dbConn)
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
		if post, err = postService.GetPost(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"author._id": bson.M{"$eq": me.UID}},
				{"_id": bson.M{"$eq": postUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if err = postService.TrashPost(post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func DetrashMyPost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			postService  = service.New(c, ctx, dbConn)
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
		if post, err = postService.GetPost(bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"author._id": bson.M{"$eq": me.UID}},
				{"_id": bson.M{"$eq": postUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, err)
			return
		}
		if err = postService.DetrashPost(post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func DeleteMyPost(
	maxCtxDuration time.Duration,
	dbConn *mongo.Database,
) (handler gin.HandlerFunc) {
	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			postService  = service.New(c, ctx, dbConn)
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
		if post, postContent, err = postService.GetPostWithContent(bson.M{
			"$and": []bson.M{
				{"author._id": bson.M{"$eq": me.UID}},
				{"_id": bson.M{"$eq": postUid}}},
		}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, err)
			return
		}
		if err = postService.DeletePost(post, postContent); err != nil {
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

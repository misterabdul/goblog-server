package posts

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/http/controllers/helpers"
	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

func GetMyPost(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me          *models.UserModel
			post        *models.PostModel
			postContent *models.PostContentModel
			postId      primitive.ObjectID
			postIdQuery = c.Param("post")
			err         error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, postContent, err = repositories.GetPostWithContent(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"_id": postId},
			}}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if post.Author.Username != me.Username {
			responses.UnauthorizedAction(c, errors.New("you are not the author of the post"))
		}

		responses.MyPost(c, post, postContent)
	}
}

func GetMyPosts(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me         *models.UserModel
			posts      []*models.PostModel
			trashQuery interface{} = primitive.Null{}
			err        error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if trashParam := c.DefaultQuery("trash", "false"); trashParam == "true" {
			trashQuery = bson.M{"$ne": primitive.Null{}}
		}
		if posts, err = repositories.GetPosts(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deletedat": trashQuery},
				{"author.username": me.Username},
			}},
			helpers.GetShowQuery(c),
			helpers.GetOrderQuery(c),
			helpers.GetAscQuery(c)); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(posts) == 0 {
			responses.NoContent(c)
		}

		responses.MyPosts(c, posts)
	}
}

func CreatePost(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me          *models.UserModel
			post        *models.PostModel
			postContent *models.PostContentModel
			form        *forms.CreatePostForm
			err         error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if form, err = requests.GetCreatePostForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		post, postContent = forms.CreatePostModel(form, me)
		if err = repositories.CreatePost(ctx, dbConn, post, postContent); err != nil {
			var writeErr mongo.WriteException
			if errors.As(err, &writeErr) {
				responses.FormIncorrect(c, err)
				return
			}
			responses.InternalServerError(c, err)
			return
		}

		responses.MyPost(c, post, postContent)
	}
}

func PublishMyPost(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me          *models.UserModel
			post        *models.PostModel
			postId      primitive.ObjectID
			postIdQuery = c.Param("post")
			err         error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, err = repositories.GetPost(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": postId},
			}}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if post.Author.Username != me.Username {
			responses.UnauthorizedAction(c, errors.New("you are not the author of the post"))
			return
		}
		if post.PublishedAt != nil {
			responses.NoContent(c)
			return
		}
		if err = repositories.PublishPost(ctx, dbConn, post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func DepublishMyPost(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me          *models.UserModel
			post        *models.PostModel
			postId      primitive.ObjectID
			postIdQuery = c.Param("post")
			err         error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, err = repositories.GetPost(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": postId},
			}}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if post.Author.Username != me.Username {
			responses.UnauthorizedAction(c, errors.New("you are not the author of the post"))
			return
		}
		if post.PublishedAt == nil {
			responses.NoContent(c)
			return
		}
		if err = repositories.DepublishPost(ctx, dbConn, post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func UpdateMyPost(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me          *models.UserModel
			post        *models.PostModel
			postContent *models.PostContentModel
			postId      primitive.ObjectID
			postIdQuery = c.Param("post")
			form        *forms.UpdatePostForm
			err         error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, postContent, err = repositories.GetPostWithContent(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": postId},
			}}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if post.Author.Username != me.Username {
			responses.UnauthorizedAction(c, errors.New("you are not the author of the post"))
			return
		}
		if form, err = requests.GetUpdatePostForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		post, postContent = forms.UpdatePostModel(form, post, postContent)
		if err = repositories.UpdatePost(ctx, dbConn, post, postContent); err != nil {
			var writeErr mongo.WriteException
			if errors.As(err, &writeErr) {
				responses.FormIncorrect(c, err)
				return
			}
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func TrashMyPost(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me          *models.UserModel
			post        *models.PostModel
			postId      primitive.ObjectID
			postIdQuery = c.Param("post")
			err         error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, err = repositories.GetPost(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": postId},
			}}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if post.Author.Username != me.Username {
			responses.UnauthorizedAction(c, errors.New("your are not the author of the post"))
			return
		}
		if err = repositories.TrashPost(ctx, dbConn, post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func DetrashMyPost(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me          *models.UserModel
			post        *models.PostModel
			postId      primitive.ObjectID
			postIdQuery = c.Param("post")
			err         error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, err = repositories.GetPost(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": postId},
			}}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, err)
			return
		}
		if post.Author.Username != me.Username {
			responses.UnauthorizedAction(c, errors.New("you are not the author of the post"))
			return
		}
		if err = repositories.DetrashPost(ctx, dbConn, post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func DeleteMyPost(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me          *models.UserModel
			post        *models.PostModel
			postContent *models.PostContentModel
			postId      primitive.ObjectID
			postIdQuery = c.Param("post")
			err         error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if post, postContent, err = repositories.GetPostWithContent(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": postId},
			}}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, err)
			return
		}
		if post.Author.Username != me.Username {
			responses.UnauthorizedAction(c, errors.New("you are not the author of the post"))
			return
		}
		if err = repositories.DeletePost(ctx, dbConn, post, postContent); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func GetMyPostComment(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me             *models.UserModel
			comment        *models.CommentModel
			post           *models.PostModel
			commentId      primitive.ObjectID
			commentIdQuery = c.Param("comment")
			err            error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if commentId, err = primitive.ObjectIDFromHex(commentIdQuery); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if comment, err = repositories.GetComment(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": commentId},
			}}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if comment == nil {
			responses.NotFound(c, errors.New("comment not found"))
			return
		}
		if post, err = repositories.GetPost(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": comment.PostUid},
			}}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if post.Author.Username != me.Username {
			responses.UnauthorizedAction(c, errors.New("you are not the author of the post"))
			return
		}

		responses.AuthorizedComment(c, comment)
	}
}

func GetMyPostComments(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me          *models.UserModel
			post        *models.PostModel
			comments    []*models.CommentModel
			postId      primitive.ObjectID
			postIdQuery             = c.Param("post")
			trashQuery  interface{} = primitive.Null{}
			err         error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if postId, err = primitive.ObjectIDFromHex(postIdQuery); err != nil {
			responses.IncorrectPostId(c, err)
			return
		}
		if trashParam := c.DefaultQuery("trash", "false"); trashParam == "true" {
			trashQuery = bson.M{"$ne": primitive.Null{}}
		}
		if post, err = repositories.GetPost(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": postId},
			}}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if post.Author.Username != me.Username {
			responses.UnauthorizedAction(c, errors.New("you are not the author of the post"))
			return
		}
		if comments, err = repositories.GetComments(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deletedat": trashQuery},
				{"postuid": post.UID},
			}},
			helpers.GetShowQuery(c),
			helpers.GetOrderQuery(c),
			helpers.GetAscQuery(c)); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(comments) == 0 {
			responses.NoContent(c)
		}

		responses.AuthorizedComments(c, comments)
	}
}

func TrashMyPostComment(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me             *models.UserModel
			comment        *models.CommentModel
			post           *models.PostModel
			commentId      primitive.ObjectID
			commentIdQuery = c.Param("comment")
			err            error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if commentId, err = primitive.ObjectIDFromHex(commentIdQuery); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if comment, err = repositories.GetComment(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deleted": primitive.Null{}},
				{"_id": commentId},
			}}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if comment == nil {
			responses.NotFound(c, errors.New("comment not found"))
			return
		}
		if post, err = repositories.GetPost(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deleted": primitive.Null{}},
				{"_id": comment.PostUid},
			}}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if post.Author.Username != me.Username {
			responses.UnauthorizedAction(c, errors.New("you are not the author of the post"))
			return
		}
		if err = repositories.TrashComment(ctx, dbConn, comment); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func DetrashMyPostComment(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me             *models.UserModel
			comment        *models.CommentModel
			post           *models.PostModel
			commentId      primitive.ObjectID
			commentIdQuery = c.Param("comment")
			err            error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if commentId, err = primitive.ObjectIDFromHex(commentIdQuery); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if comment, err = repositories.GetComment(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": commentId},
			}}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if comment == nil {
			responses.NotFound(c, errors.New("comment not found"))
			return
		}
		if post, err = repositories.GetPost(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": comment.PostUid},
			}}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if post.Author.Username != me.Username {
			responses.UnauthorizedAction(c, errors.New("you are not the author of the post"))
			return
		}
		if err = repositories.DetrashComment(ctx, dbConn, comment); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

func DeleteMyPostComment(maxCtxDuration time.Duration, dbConn *mongo.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			me             *models.UserModel
			comment        *models.CommentModel
			post           *models.PostModel
			commentId      primitive.ObjectID
			commentIdQuery = c.Param("comment")
			err            error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		if commentId, err = primitive.ObjectIDFromHex(commentIdQuery); err != nil {
			responses.IncorrectCommentId(c, err)
			return
		}
		if comment, err = repositories.GetComment(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deletedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": commentId},
			}}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if comment == nil {
			responses.NotFound(c, errors.New("comment not found"))
			return
		}
		if post, err = repositories.GetPost(ctx, dbConn,
			bson.M{"$and": []bson.M{
				{"deletedat": primitive.Null{}},
				{"_id": comment.PostUid},
			}}); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if post.Author.Username != me.Username {
			responses.UnauthorizedAction(c, errors.New("you are not the author of the post"))
			return
		}
		if err = repositories.DeleteComment(ctx, dbConn, comment); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.NoContent(c)
	}
}

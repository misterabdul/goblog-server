package posts

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/internal/http/controllers/helpers"
	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
)

func GetPublicPostComment(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn    *mongo.Database
			comment   *models.CommentModel
			commentId primitive.ObjectID
			err       error
		)

		commentIdQuery := c.Param("comment")
		if commentId, err = primitive.ObjectIDFromHex(commentIdQuery); err != nil {
			responses.NotFound(c, err)
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if comment, err = repositories.GetComment(ctx, dbConn,
			bson.M{"$and": []interface{}{
				bson.M{"deletedat": primitive.Null{}},
				bson.M{"_id": commentId},
			}}); err != nil {
			responses.NotFound(c, err)
			return
		}
		if _, err = repositories.GetPost(ctx, dbConn,
			bson.M{"$and": []interface{}{
				bson.M{"deletedat": primitive.Null{}},
				bson.M{"publishedat": bson.M{"$ne": primitive.Null{}}},
				bson.M{"_id": comment.PostUid},
			}}); err != nil {
			responses.NotFound(c, err)
			return
		}

		responses.PublicComment(c, comment)
	}
}

func GetPublicPostComments(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn   *mongo.Database
			comments []*models.CommentModel
			postId   primitive.ObjectID
			err      error
		)

		postQuery := c.Param("post")
		if postId, err = primitive.ObjectIDFromHex(postQuery); err != nil {
			postId = primitive.ObjectID{}
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if _, err = repositories.GetPost(ctx, dbConn,
			bson.M{"$and": []interface{}{
				bson.M{"deletedat": primitive.Null{}},
				bson.M{"publishedat": bson.M{"$ne": primitive.Null{}}},
				bson.M{"$or": []interface{}{
					bson.M{"_id": postId},
					bson.M{"slug": postQuery},
				}},
			}}); err != nil {
			responses.NotFound(c, err)
			return
		}
		if comments, err = repositories.GetComments(ctx, dbConn,
			bson.M{"$and": []interface{}{
				bson.M{"deletedat": primitive.Null{}},
				bson.M{"$or": []interface{}{
					bson.M{"postslug": postQuery},
					bson.M{"postuid": postId},
				}},
			}},
			helpers.GetShowQuery(c),
			helpers.GetOrderQuery(c),
			helpers.GetAscQuery(c)); err != nil {
			responses.NotFound(c, err)
		}

		responses.PublicComments(c, comments)
	}
}

func CreatePublicPostComment(maxCtxDuration time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), maxCtxDuration)
		defer cancel()

		var (
			dbConn  *mongo.Database
			comment *models.CommentModel
			form    *forms.CreateCommentForm
			err     error
		)

		if form, err = requests.GetCreateCommentForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if comment, err = forms.CreateCommentModel(form); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer dbConn.Client().Disconnect(ctx)

		if _, err = repositories.GetPost(ctx, dbConn,
			bson.M{"$and": []interface{}{
				bson.M{"deletedat": primitive.Null{}},
				bson.M{"publishedat": bson.M{"$ne": primitive.Null{}}},
				bson.M{"_id": comment.PostUid},
			}}); err != nil {
			responses.NotFound(c, err)
			return
		}
		if err = repositories.CreateComment(ctx, dbConn, comment); err != nil {
			var writeErr mongo.WriteException
			if errors.As(err, &writeErr) {
				responses.FormIncorrect(c, err)
				return
			}
			responses.InternalServerError(c, err)
			return
		}

		responses.PublicComment(c, comment)
	}
}

package comments

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/http/forms"
	"github.com/misterabdul/goblog-server/internal/http/requests"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	"github.com/misterabdul/goblog-server/internal/service"
)

// @Tags        Comment (Public)
// @Summary     Get Public Comment
// @Description Get a comment that available publicly.
// @Router      /v1/comment/{uid} [get]
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Comment's UID"
// @Success     200 {object} object{data=object{uid=string,postUid=string,parentCommentUid=string,email=string,name=string,content=string,replyCount=int,createdAt=time}}
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func GetPublicComment(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel     = context.WithTimeout(context.Background(), maxCtxDuration)
			comment         *models.CommentModel
			post            *models.PostModel
			commentUid      primitive.ObjectID
			commentUidParam = c.Param("comment")
			err             error
		)

		defer cancel()
		if commentUid, err = primitive.ObjectIDFromHex(commentUidParam); err != nil {
			responses.NotFound(c, err)
			return
		}
		if comment, err = svc.Comment.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"_id": bson.M{"$eq": commentUid}}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if comment == nil {
			responses.NotFound(c, errors.New("comment not found"))
			return
		}
		if post, err = svc.Post.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": bson.M{"$eq": comment.PostUid}}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
		}

		responses.PublicComment(c, comment)
	}
}

// @Tags        Comment (Public)
// @Summary     Get Public Post's Comments
// @Description Get public post's comments that available publicly.
// @Router      /v1/post/{uid}/comments [get]
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Post's UID"
// @Success     200 {object} object{data=[]object{uid=string,postUid=string,parentCommentUid=string,email=string,name=string,content=string,replyCount=int,createdAt=time}}
// @Failure     204
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func GetPublicPostComments(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			comments    []*models.CommentModel
			post        *models.PostModel
			postUid     interface{}
			postParam   = c.Param("post")
			err         error
		)

		defer cancel()
		if postUid, err = primitive.ObjectIDFromHex(postParam); err != nil {
			postUid = nil
		}
		if post, err = svc.Post.GetOne(ctx, bson.M{
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
		if comments, err = svc.Comment.GetMany(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"parentcommentuid": bson.M{"$eq": primitive.Null{}}},
				{"postuid": bson.M{"$eq": post.UID}}}},
			internalGin.GetFindOptions(c),
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(comments) == 0 {
			responses.NoContent(c)
			return
		}

		responses.PublicComments(c, comments)
	}
}

// @Tags        Comment (Public)
// @Summary     Get Public Comment's Replies
// @Description Get public comment's replies that available publicly.
// @Router      /v1/comment/{uid}/replies [get]
// @Produce     application/json
// @Produce     application/msgpack
// @Param       uid path     string true "Comment's UID"
// @Success     200 {object} object{data=[]object{uid=string,postUid=string,parentCommentUid=string,email=string,name=string,content=string,replyCount=int,createdAt=time}}
// @Failure     204
// @Failure     404 {object} object{message=string}
// @Failure     500 {object} object{message=string}
func GetPublicCommentReplies(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel  = context.WithTimeout(context.Background(), maxCtxDuration)
			replies      []*models.CommentModel
			comment      *models.CommentModel
			post         *models.PostModel
			commentUid   interface{}
			commentParam = c.Param("comment")
			err          error
		)

		defer cancel()
		if commentUid, err = primitive.ObjectIDFromHex(commentParam); err != nil {
			responses.NotFound(c, errors.New("incorrent comment id format"))
			return
		}
		if comment, err = svc.Comment.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"_id": bson.M{"$eq": commentUid}}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if comment == nil {
			responses.NotFound(c, errors.New("comment not found"))
			return
		}
		if post, err = svc.Post.GetOne(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"publishedat": bson.M{"$ne": primitive.Null{}}},
				{"_id": bson.M{"$eq": comment.PostUid}}}},
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if post == nil {
			responses.NotFound(c, errors.New("post not found"))
			return
		}
		if replies, err = svc.Comment.GetMany(ctx, bson.M{
			"$and": []bson.M{
				{"deletedat": bson.M{"$eq": primitive.Null{}}},
				{"postuid": bson.M{"$eq": post.UID}},
				{"parentcommentuid": bson.M{"$eq": comment.UID}}}},
			internalGin.GetFindOptions(c),
		); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if len(replies) == 0 {
			responses.NoContent(c)
			return
		}

		responses.PublicComments(c, replies)
	}
}

// @Tags        Comment (Public)
// @Summary     Create Public Post's Comment
// @Description Create a comment for a post that available publicly.
// @Router      /v1/comment [post]
// @Accept      application/json
// @Accept      application/msgpack
// @Produce     application/json
// @Produce     application/msgpack
// @Param       form body     object{postUid=string,email=string,name=string,content=string} true "Create comment form"
// @Success     200  {object} object{data=object{uid=string,postUid=string,parentCommentUid=string,email=string,name=string,content=string,replyCount=int,createdAt=time}}
// @Failure     404  {object} object{message=string}
// @Failure     422  {object} object{message=string}
// @Failure     500  {object} object{message=string}
func CreatePublicPostComment(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			comment     *models.CommentModel
			post        *models.PostModel
			form        *forms.CreateCommentForm
			err         error
		)

		defer cancel()
		if form, err = requests.GetCreateCommentForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if post, err = form.Validate(svc, ctx); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if comment, err = form.ToCommentModel(); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = svc.Comment.SaveOne(ctx, comment, post); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.PublicComment(c, comment)
	}
}

// @Tags        Comment (Public)
// @Summary     Create Public Comment's Reply
// @Description Create a reply for a comment that available publicly.
// @Router      /v1/comment/reply [post]
// @Accept      application/json
// @Accept      application/msgpack
// @Produce     application/json
// @Produce     application/msgpack
// @Param       form body     object{parentCommentUid=string,email=string,name=string,content=string} true "Create comment form"
// @Success     200  {object} object{data=object{uid=string,postUid=string,parentCommentUid=string,email=string,name=string,content=string,replyCount=int,createdAt=time}}
// @Failure     404  {object} object{message=string}
// @Failure     422  {object} object{message=string}
// @Failure     500  {object} object{message=string}
func CreatePublicCommentReply(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			ctx, cancel = context.WithTimeout(context.Background(), maxCtxDuration)
			reply       *models.CommentModel
			comment     *models.CommentModel
			form        *forms.CreateCommentReplyForm
			err         error
		)

		defer cancel()
		if form, err = requests.GetCreateCommentReplyForm(c); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if comment, err = form.Validate(svc, ctx); err != nil {
			responses.FormIncorrect(c, err)
			return
		}
		if reply, err = form.ToCommentReplyModel(); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		if err = svc.Comment.SaveOneReply(ctx, reply, comment); err != nil {
			responses.InternalServerError(c, err)
			return
		}

		responses.PublicComment(c, reply)
	}
}

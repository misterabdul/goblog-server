package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/models"
)

func PublicComment(c *gin.Context, comment *models.CommentModel) {
	data := extractPublicCommentData(comment)
	Basic(c, http.StatusOK, data)
}

func AuthorizedComment(c *gin.Context, comment *models.CommentModel) {
	data := extractAuthorizedCommentData(comment)
	Basic(c, http.StatusOK, data)
}

func PublicComments(c *gin.Context, comments []*models.CommentModel) {
	var data []gin.H

	for _, comment := range comments {
		data = append(data, extractPublicCommentData(comment))
	}
	Basic(c, http.StatusOK, gin.H{"data": data})
}

func AuthorizedComments(c *gin.Context, comments []*models.CommentModel) {
	var data []gin.H

	for _, comment := range comments {
		data = append(data, extractAuthorizedCommentData(comment))
	}
	Basic(c, http.StatusOK, gin.H{"data": data})
}

func IncorrectCommentId(c *gin.Context, err error) {
	Basic(c, http.StatusBadRequest, gin.H{
		"message": "incorrent comment id format"})
}

func extractPublicCommentData(comment *models.CommentModel) (extracted gin.H) {
	return gin.H{
		"uid":              comment.UID.Hex(),
		"postUid":          comment.PostUid,
		"parentCommentUid": comment.ParentCommentUid,
		"email":            comment.Email,
		"name":             comment.Name,
		"content":          comment.Content,
		"replyCount":       comment.ReplyCount,
		"createdAt":        comment.CreatedAt}
}

func extractAuthorizedCommentData(comment *models.CommentModel) (extracted gin.H) {
	return gin.H{
		"uid":              comment.UID.Hex(),
		"postUid":          comment.PostUid,
		"postAuthorUid":    comment.PostAuthorUid,
		"parentCommentUid": comment.ParentCommentUid,
		"email":            comment.Email,
		"name":             comment.Name,
		"content":          comment.Content,
		"replyCount":       comment.ReplyCount,
		"createdAt":        comment.CreatedAt,
		"deletedAt":        comment.DeletedAt}
}

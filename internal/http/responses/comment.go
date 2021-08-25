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
	Basic(c, http.StatusBadRequest, gin.H{"message": "incorrent comment id format"})
}

func extractPublicCommentData(comment *models.CommentModel) gin.H {
	return gin.H{
		"uid":       comment.UID,
		"postUid":   comment.PostUid,
		"postSlug":  comment.PostSlug,
		"email":     comment.Email,
		"name":      comment.Name,
		"content":   comment.Content,
		"replies":   extractPublicCommentRepliesData(comment.Replies),
		"createdAt": comment.CreatedAt,
	}
}

func extractPublicCommentRepliesData(replies []models.CommentReplyModel) []gin.H {
	var data []gin.H

	for _, reply := range replies {
		data = append(data, gin.H{
			"email":     reply.Email,
			"name":      reply.Name,
			"content":   reply.Content,
			"createdAt": reply.CreatedAt,
		})
	}

	return data
}

func extractAuthorizedCommentData(comment *models.CommentModel) gin.H {
	return gin.H{
		"uid":       comment.UID,
		"postUid":   comment.PostUid,
		"postSlug":  comment.PostSlug,
		"email":     comment.Email,
		"name":      comment.Name,
		"content":   comment.Content,
		"replies":   extractAuthorizedCommentRepliesData(comment.Replies),
		"createdAt": comment.CreatedAt,
		"deletedAt": comment.DeletedAt,
	}
}

func extractAuthorizedCommentRepliesData(replies []models.CommentReplyModel) []gin.H {
	var data []gin.H

	for _, reply := range replies {
		data = append(data, gin.H{
			"email":     reply.Email,
			"name":      reply.Name,
			"content":   reply.Content,
			"createdAt": reply.CreatedAt,
			"deletedAt": reply.DeletedAt,
		})
	}

	return data
}

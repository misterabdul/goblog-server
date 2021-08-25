package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/models"
)

func PublicPost(c *gin.Context, post *models.PostModel) {
	data := extractPublicPostData(post)

	Basic(c, http.StatusOK, data)
}

func AuthorizedPost(c *gin.Context, post *models.PostModel) {
	data := extractAuthorizedPostData(post)

	Basic(c, http.StatusOK, data)
}

func MyPost(c *gin.Context, post *models.PostModel) {
	AuthorizedPost(c, post)
}

func PublicPosts(c *gin.Context, posts []*models.PostModel) {
	var data []gin.H
	for _, post := range posts {
		data = append(data, extractPublicPostData(post))
	}

	Basic(c, http.StatusOK, gin.H{"data": data})
}

func AuthorizedPosts(c *gin.Context, posts []*models.PostModel) {
	var data []gin.H
	for _, post := range posts {
		data = append(data, extractAuthorizedPostData(post))
	}

	Basic(c, http.StatusOK, gin.H{"data": data})
}

func MyPosts(c *gin.Context, posts []*models.PostModel) {
	AuthorizedPosts(c, posts)
}

func IncorrectPostId(c *gin.Context, err error) {
	Basic(c, http.StatusBadRequest, gin.H{"message": "incorrent post id format"})
}

func extractPublicPostData(post *models.PostModel) gin.H {
	return gin.H{
		"uid":         post.UID,
		"slug":        post.Slug,
		"title":       post.Title,
		"categories":  post.Categories,
		"tags":        post.Tags,
		"content":     post.Content,
		"author":      post.Author,
		"publishedAt": post.PublishedAt,
	}
}

func extractAuthorizedPostData(post *models.PostModel) gin.H {
	return gin.H{
		"uid":         post.UID,
		"slug":        post.Slug,
		"title":       post.Title,
		"categories":  post.Categories,
		"tags":        post.Tags,
		"content":     post.Content,
		"author":      post.Author,
		"publishedAt": post.PublishedAt,
		"createdAt":   post.CreatedAt,
		"updatedAt":   post.UpdatedAt,
		"deletedat":   post.DeletedAt,
	}
}

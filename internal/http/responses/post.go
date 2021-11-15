package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/models"
)

func PublicPost(
	c *gin.Context,
	post *models.PostModel,
	postContent *models.PostContentModel,
) {
	data := extractPublicPostData(post, postContent)
	Basic(c, http.StatusOK, data)
}

func AuthorizedPost(
	c *gin.Context,
	post *models.PostModel,
	postContent *models.PostContentModel,
) {
	data := extractAuthorizedPostData(post, postContent)
	Basic(c, http.StatusOK, data)
}

func MyPost(
	c *gin.Context,
	post *models.PostModel,
	postContent *models.PostContentModel,
) {
	AuthorizedPost(c, post, postContent)
}

func PublicPosts(c *gin.Context, posts []*models.PostModel) {
	var data []gin.H

	for _, post := range posts {
		data = append(data, extractPublicPostData(post, nil))
	}
	Basic(c, http.StatusOK, gin.H{"data": data})
}

func AuthorizedPosts(c *gin.Context, posts []*models.PostModel) {
	var data []gin.H

	for _, post := range posts {
		data = append(data, extractAuthorizedPostData(post, nil))
	}
	Basic(c, http.StatusOK, gin.H{"data": data})
}

func MyPosts(c *gin.Context, posts []*models.PostModel) {
	AuthorizedPosts(c, posts)
}

func IncorrectPostId(c *gin.Context, err error) {
	Basic(c, http.StatusBadRequest, gin.H{
		"message": "incorrect post id format"})
}

func extractPublicPostData(
	post *models.PostModel,
	postContent *models.PostContentModel,
) (extracted gin.H) {
	if postContent == nil || postContent.UID != post.UID {
		return gin.H{
			"uid":                post.UID,
			"slug":               post.Slug,
			"title":              post.Title,
			"featuringImagePath": post.FeaturingImagePath,
			"description":        post.Description,
			"categories":         post.Categories,
			"tags":               post.Tags,
			"author":             post.Author,
			"publishedAt":        post.PublishedAt}
	}
	return gin.H{
		"uid":                post.UID,
		"slug":               post.Slug,
		"title":              post.Title,
		"featuringImagePath": post.FeaturingImagePath,
		"description":        post.Description,
		"categories":         post.Categories,
		"tags":               post.Tags,
		"content":            postContent.Content,
		"author":             post.Author,
		"publishedAt":        post.PublishedAt}
}

func extractAuthorizedPostData(
	post *models.PostModel,
	postContent *models.PostContentModel,
) (extracted gin.H) {
	if postContent == nil || postContent.UID != post.UID {
		return gin.H{
			"uid":                post.UID,
			"slug":               post.Slug,
			"title":              post.Title,
			"featuringImagePath": post.FeaturingImagePath,
			"description":        post.Description,
			"categories":         post.Categories,
			"tags":               post.Tags,
			"author":             post.Author,
			"publishedAt":        post.PublishedAt,
			"createdAt":          post.CreatedAt,
			"updatedAt":          post.UpdatedAt,
			"deletedat":          post.DeletedAt}
	}
	return gin.H{
		"uid":                post.UID,
		"slug":               post.Slug,
		"title":              post.Title,
		"featuringImagePath": post.FeaturingImagePath,
		"description":        post.Description,
		"categories":         post.Categories,
		"tags":               post.Tags,
		"content":            postContent.Content,
		"author":             post.Author,
		"publishedAt":        post.PublishedAt,
		"createdAt":          post.CreatedAt,
		"updatedAt":          post.UpdatedAt,
		"deletedat":          post.DeletedAt}
}

package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/database/models"
)

func PublicPage(
	c *gin.Context,
	page *models.PageModel,
	pageContent *models.PageContentModel,
) {
	data := extractPublicPageData(page, pageContent)
	Basic(c, http.StatusOK, gin.H{"data": data})
}

func AuthorizedPage(
	c *gin.Context,
	page *models.PageModel,
	pageContent *models.PageContentModel,
) {
	data := extractAuthorizedPageData(page, pageContent)
	Basic(c, http.StatusOK, gin.H{"data": data})
}

func MyPage(
	c *gin.Context,
	page *models.PageModel,
	pageContent *models.PageContentModel,
) {
	AuthorizedPage(c, page, pageContent)
}

func PublicPages(c *gin.Context, pages []*models.PageModel) {
	var data []gin.H

	for _, page := range pages {
		data = append(data, extractPublicPageData(page, nil))
	}
	Basic(c, http.StatusOK, gin.H{"data": data})
}

func AuthorizedPages(c *gin.Context, pages []*models.PageModel) {
	var data []gin.H

	for _, page := range pages {
		data = append(data, extractAuthorizedPageData(page, nil))
	}
	Basic(c, http.StatusOK, gin.H{"data": data})
}

func MyPages(c *gin.Context, pages []*models.PageModel) {
	AuthorizedPages(c, pages)
}

func IncorrectPageId(c *gin.Context, err error) {
	Basic(c, http.StatusBadRequest, gin.H{
		"message": "incorrect page id format"})
}

func extractPublicPageData(
	page *models.PageModel,
	pageContent *models.PageContentModel,
) (extracted gin.H) {
	if pageContent == nil || pageContent.UID != page.UID {
		return gin.H{
			"uid":         page.UID.Hex(),
			"slug":        page.Slug,
			"title":       page.Title,
			"publishedAt": page.PublishedAt}
	}
	return gin.H{
		"uid":         page.UID.Hex(),
		"slug":        page.Slug,
		"title":       page.Title,
		"content":     pageContent.Content,
		"publishedAt": page.PublishedAt}
}

func extractAuthorizedPageData(
	page *models.PageModel,
	pageContent *models.PageContentModel,
) (extracted gin.H) {
	if pageContent == nil || pageContent.UID != page.UID {
		return gin.H{
			"uid":         page.UID.Hex(),
			"slug":        page.Slug,
			"title":       page.Title,
			"author":      extractCommonAuthorData(page.Author),
			"publishedAt": page.PublishedAt,
			"createdAt":   page.CreatedAt,
			"updatedAt":   page.UpdatedAt,
			"deletedat":   page.DeletedAt}
	}
	return gin.H{
		"uid":         page.UID.Hex(),
		"slug":        page.Slug,
		"title":       page.Title,
		"content":     pageContent.Content,
		"author":      extractCommonAuthorData(page.Author),
		"publishedAt": page.PublishedAt,
		"createdAt":   page.CreatedAt,
		"updatedAt":   page.UpdatedAt,
		"deletedat":   page.DeletedAt}
}

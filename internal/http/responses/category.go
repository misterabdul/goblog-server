package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/models"
)

func PublicCategory(c *gin.Context, category *models.CategoryModel) {
	data := extractPublicCategoryData(category)
	Basic(c, http.StatusOK, gin.H{"data": data})
}

func AuthorizedCategory(c *gin.Context, category *models.CategoryModel) {
	data := extractAuthorizedCategoryData(category)
	Basic(c, http.StatusOK, gin.H{"data": data})
}

func PublicCategories(c *gin.Context, categories []*models.CategoryModel) {
	var data []gin.H

	for _, category := range categories {
		data = append(data, extractPublicCategoryData(category))
	}
	Basic(c, http.StatusOK, gin.H{"data": data})
}

func AuthorizedCategories(c *gin.Context, categories []*models.CategoryModel) {
	var data []gin.H

	for _, category := range categories {
		data = append(data, extractAuthorizedCategoryData(category))
	}
	Basic(c, http.StatusOK, gin.H{"data": data})
}

func IncorrectCategoryId(c *gin.Context, err error) {
	Basic(c, http.StatusBadRequest, gin.H{
		"message": "incorrent post id format"})
}

func extractPublicCategoryData(category *models.CategoryModel) (extracted gin.H) {
	return gin.H{
		"uid":  category.UID.Hex(),
		"slug": category.Slug,
		"name": category.Name}
}

func extractAuthorizedCategoryData(category *models.CategoryModel) (extracted gin.H) {
	return gin.H{
		"uid":       category.UID.Hex(),
		"slug":      category.Slug,
		"name":      category.Name,
		"createdAt": category.CreatedAt,
		"updatedAt": category.UpdatedAt,
		"deletedat": category.DeletedAt}
}

func extractPostCategoryData(categories []models.CategoryCommonModel) (extracted []gin.H) {
	for _, category := range categories {
		extracted = append(extracted, gin.H{
			"uid":  category.UID,
			"slug": category.Slug,
			"name": category.Name,
		})
	}

	return extracted
}

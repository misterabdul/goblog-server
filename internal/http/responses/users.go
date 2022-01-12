package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/models"
)

func PublicUser(c *gin.Context, user *models.UserModel) {
	data := extractPublicUserData(user)
	Basic(c, http.StatusOK, gin.H{"data": data})
}

func AuthorizedUser(c *gin.Context, user *models.UserModel) {
	data := extractAuthorizedUserData(user)
	Basic(c, http.StatusOK, gin.H{"data": data})
}

func Me(c *gin.Context, user *models.UserModel) {
	data := extractAuthorizedUserData(user)
	Basic(c, http.StatusOK, gin.H{"data": data})
}

func PublicUsers(c *gin.Context, users []*models.UserModel) {
	var data []gin.H

	for _, user := range users {
		data = append(data, extractPublicUserData(user))
	}
	Basic(c, http.StatusOK, gin.H{"data": data})
}

func AuthorizedUsers(c *gin.Context, users []*models.UserModel) {
	var data []gin.H

	for _, user := range users {
		data = append(data, extractAuthorizedUserData(user))
	}
	Basic(c, http.StatusOK, gin.H{"data": data})
}

func IncorrectUserId(c *gin.Context, err error) {
	Basic(c, http.StatusBadRequest, gin.H{
		"message": "incorrect user id format"})
}

func extractPublicUserData(user *models.UserModel) gin.H {
	return gin.H{
		"uid":       user.UID.Hex(),
		"username":  user.Username,
		"email":     user.Email,
		"firstName": user.FirstName,
		"lastName":  user.LastName}
}

func extractAuthorizedUserData(user *models.UserModel) gin.H {
	return gin.H{
		"uid":       user.UID.Hex(),
		"username":  user.Username,
		"email":     user.Email,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"roles":     user.Roles,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt}
}

func extractPostAuthorData(user models.UserCommonModel) gin.H {
	return gin.H{
		"username":  user.Username,
		"email":     user.Email,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
	}
}

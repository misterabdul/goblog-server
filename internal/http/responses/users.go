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

func extractPublicUserData(user *models.UserModel) gin.H {
	return gin.H{
		"uid":       user.UID,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"username":  user.Username,
	}
}

func extractAuthorizedUserData(user *models.UserModel) gin.H {
	return gin.H{
		"uid":       user.UID,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"username":  user.Username,
		"email":     user.Email,
		"roles":     user.Roles,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	}
}

package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/misterabdul/goblog-server/internal/models"
)

func User(c *gin.Context, user *models.UserModel) {
	data := extractUserData(user)

	Basic(c, http.StatusOK, gin.H{"data": data})
}

func Me(c *gin.Context, user *models.UserModel) {
	User(c, user)
}

func Users(c *gin.Context, users []*models.UserModel) {
	var data []gin.H
	for _, user := range users {
		data = append(data, extractUserData(user))
	}

	Basic(c, http.StatusOK, gin.H{"data": data})
}

func extractUserData(user *models.UserModel) gin.H {
	return gin.H{
		"uid":       user.UID,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"username":  user.Username,
		"email":     user.Email,
	}
}

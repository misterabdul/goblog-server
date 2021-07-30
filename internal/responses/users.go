package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/misterabdul/goblog-server/internal/models"
)

func User(c *gin.Context, user *models.UserModel) {
	data := gin.H{
		"uid":       user.UID,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
	}

	Basic(c, http.StatusOK, gin.H{"data": data})
}

func Users(c *gin.Context, users []*models.UserModel) {
	var data []interface{}
	for _, user := range users {
		data = append(data, gin.H{
			"uid":       user.UID,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
		})
	}

	Basic(c, http.StatusOK, gin.H{"data": data})
}

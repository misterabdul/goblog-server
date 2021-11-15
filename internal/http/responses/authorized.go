package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UnauthorizedAction(c *gin.Context, err error) {
	Basic(c, http.StatusUnauthorized, gin.H{
		"message": "unauthorized action"})
}

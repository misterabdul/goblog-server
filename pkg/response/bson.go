package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// BSON http response for gin.
func BSON(c *gin.Context, code int, obj interface{}) {
	bsonData, err := bson.Marshal(obj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.Data(code, "application/bson", bsonData)
}

package responses

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"

	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
)

func ResourceStats(c *gin.Context, count int64) {
	var (
		show       = internalGin.GetShowQuery(c)
		page       = internalGin.GetPageQuery(c)
		totalPages = math.Ceil(float64(count) / float64(*show))
	)

	Basic(c, http.StatusOK, gin.H{"data": gin.H{
		"currentPage":  *page,
		"totalPages":   totalPages,
		"itemsPerPage": *show,
		"totalItems":   count}})
}

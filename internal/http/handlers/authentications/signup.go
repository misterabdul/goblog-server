package authentications

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/http/responses"
	"github.com/misterabdul/goblog-server/internal/service"
)

func SignUp(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		responses.NotImplemented(c, nil)
	}
}

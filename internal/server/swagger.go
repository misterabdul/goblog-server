package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	swaggerGin "github.com/swaggo/gin-swagger"

	swaggerDocs "github.com/misterabdul/goblog-server/api/swagger"
)

func InitSwagger(
	server *gin.Engine,
) {
	var serverEnv *serverRelatedEnv = getServerRelatedEnv()

	switch serverEnv.Mode {
	case 1:
		fallthrough
	case 2:
		swaggerDocs.SwaggerInfo.Host = ReadAddressFromEnv()
		server.GET("/swagger/*any", swaggerGin.WrapHandler(swaggerFiles.Handler))
	}
}

package server

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type serverRelatedEnv struct {
	Mode int
}

// Get the server engine.
func GetServer() (serverInstance *gin.Engine) {
	var (
		serverEnv  *serverRelatedEnv = getServerRelatedEnv()
		corsConfig cors.Config
	)

	switch serverEnv.Mode {
	default:
		fallthrough
	case 0:
		gin.SetMode(gin.ReleaseMode)
		serverInstance = gin.New()
		serverInstance.Use(gin.Recovery())
	case 1:
		gin.SetMode(gin.TestMode)
		serverInstance = gin.New()
	case 2:
		gin.SetMode(gin.DebugMode)
		serverInstance = gin.Default()
		corsConfig = cors.DefaultConfig()
		corsConfig.AllowAllOrigins = true
		corsConfig.AllowCredentials = true
		serverInstance.Use(cors.New(corsConfig))
	}
	serverInstance.Use(gzip.Gzip(gzip.DefaultCompression))

	return serverInstance
}

func getServerRelatedEnv() (envs *serverRelatedEnv) {
	_envs := serverRelatedEnv{
		Mode: 0,
	}

	if mode, ok := os.LookupEnv("APP_MODE"); ok {
		switch {
		default:
			fallthrough
		case mode == "release":
			_envs.Mode = 0
		case mode == "test":
			_envs.Mode = 1
		case mode == "debug":
			_envs.Mode = 2
		}
	}

	return &_envs
}

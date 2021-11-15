package server

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type serverRelatedEnv struct {
	Mode int
}

// Get the server engine.
func GetServer(dbConn *mongo.Database) (serverInstance *gin.Engine) {
	serverEnv := getServerRelatedEnv()
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
		serverInstance.Use(cors.Default())
	}
	serverInstance.Use(gzip.Gzip(gzip.DefaultCompression))
	initRoute(serverInstance, dbConn)

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

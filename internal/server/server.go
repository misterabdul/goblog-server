package server

import (
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type serverRelatedEnv struct {
	Mode int
}

// Get the server engine.
func GetServer(dbConn *mongo.Database) *gin.Engine {
	var (
		ginEngine *gin.Engine
	)

	serverEnv := getServerRelatedEnv()
	switch serverEnv.Mode {
	default:
		fallthrough
	case 0:
		gin.SetMode(gin.ReleaseMode)
		ginEngine = gin.New()
		ginEngine.Use(gin.Recovery())
	case 1:
		gin.SetMode(gin.TestMode)
		ginEngine = gin.New()
	case 2:
		gin.SetMode(gin.DebugMode)
		ginEngine = gin.Default()
	}

	initRoute(ginEngine, dbConn)

	return ginEngine
}

func getServerRelatedEnv() *serverRelatedEnv {
	env := serverRelatedEnv{
		Mode: 0,
	}

	if mode, ok := os.LookupEnv("APP_MODE"); ok {
		switch {
		default:
			fallthrough
		case mode == "release":
			env.Mode = 0
		case mode == "test":
			env.Mode = 1
		case mode == "debug":
			env.Mode = 2
		}
	}

	return &env
}

package server

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type serverRelatedEnv struct {
	Mode           int
	TrustedProxies *[]string
	UseCors        bool
	CorsMiddleware *gin.HandlerFunc
}

// Get the server engine.
func GetServer() (serverInstance *gin.Engine) {
	var serverEnv *serverRelatedEnv = getServerRelatedEnv()

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
	}
	if serverEnv.TrustedProxies != nil {
		serverInstance.SetTrustedProxies(*serverEnv.TrustedProxies)
	}
	if serverEnv.UseCors {
		serverInstance.Use(*serverEnv.CorsMiddleware)
	}
	serverInstance.Use(gzip.Gzip(gzip.DefaultCompression))

	return serverInstance
}

func getServerRelatedEnv() (envs *serverRelatedEnv) {
	var (
		_envs              serverRelatedEnv
		_trustedProxiesEnv string
		_trustedProxies    []string
		_allowedOriginsEnv string
		corsConfig         cors.Config
		corsMiddleware     gin.HandlerFunc
		ok                 bool
	)

	_envs = serverRelatedEnv{
		Mode:           0,
		TrustedProxies: nil,
		UseCors:        false,
		CorsMiddleware: nil,
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

	if _trustedProxiesEnv, ok = os.LookupEnv("TRUSTED_PROXIES"); ok {
		_trustedProxies = strings.Split(_trustedProxiesEnv, ",")
		_envs.TrustedProxies = &_trustedProxies
	}

	if _allowedOriginsEnv, ok = os.LookupEnv("CORS_ALLOWED_ORIGINS"); ok {
		corsConfig = cors.DefaultConfig()
		corsConfig.AllowAllOrigins = false
		corsConfig.AllowOrigins = strings.Split(_allowedOriginsEnv, ",")
		corsConfig.AllowCredentials = true
		corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "authorization")
		corsMiddleware = cors.New(corsConfig)
		_envs.CorsMiddleware = &corsMiddleware
		_envs.UseCors = true
	}

	return &_envs
}

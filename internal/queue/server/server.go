package server

import (
	"os"

	"github.com/hibiken/asynq"
)

type serverRelatedEnv struct {
	Mode int
}

func GetServer() *asynq.Server {
	var (
		serverEnv = getRedisServerRelatedEnv()
		logLevel  = asynq.ErrorLevel
	)

	switch serverEnv.Mode {
	default:
		fallthrough
	case 0:
		logLevel = asynq.ErrorLevel
	case 1:
		logLevel = asynq.WarnLevel
	case 2:
		logLevel = asynq.InfoLevel
	}

	return asynq.NewServer(
		ReadRedisOptsFromEnv(),
		asynq.Config{
			Concurrency: 10,
			LogLevel:    logLevel,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1}})
}

func ReadRedisOptsFromEnv() asynq.RedisClientOpt {
	var (
		host    = "localhost"
		port    = "6379"
		user    = ""
		pass    = ""
		envHost string
		envPort string
		envUser string
		envPass string
		ok      bool
	)

	if envHost, ok = os.LookupEnv("REDIS_HOST"); ok {
		host = envHost
	} else {
		host = "localhost"
	}
	if envPort, ok = os.LookupEnv("REDIS_PORT"); ok {
		port = envPort
	} else {
		port = "80"
	}
	if envUser, ok = os.LookupEnv("REDIS_USER"); ok {
		user = envUser
	} else {
		user = ""
	}
	if envPass, ok = os.LookupEnv("REDIS_PASS"); ok {
		pass = envPass
	} else {
		pass = "80"
	}

	return asynq.RedisClientOpt{
		Addr:     host + ":" + port,
		Username: user,
		Password: pass,
	}
}

func getRedisServerRelatedEnv() (envs *serverRelatedEnv) {
	var (
		_envs serverRelatedEnv
		mode  string
		ok    bool
	)

	_envs = serverRelatedEnv{
		Mode: 0}
	if mode, ok = os.LookupEnv("APP_MODE"); ok {
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
	} else {
		_envs.Mode = 0
	}

	return &_envs
}

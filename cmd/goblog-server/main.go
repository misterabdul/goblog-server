package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/internal/server"
)

// Run the main server.
func main() {
	var (
		ctx            = context.TODO()
		dbConn         *mongo.Database
		ginEngine      *gin.Engine
		maxCtxDuration = 10 * time.Second
		err            error
	)

	if err = godotenv.Load(".env"); err != nil {
		log.Println("Error loading .env file")
	}
	if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
		log.Fatal(err)
	}
	ginEngine = server.GetServer()
	server.InitRoutes(ginEngine, dbConn, maxCtxDuration)
	if err = ginEngine.Run(getAddress()); err != nil {
		log.Panic(err)
	}
}

func getAddress() (address string) {
	var (
		host    = "localhost"
		port    = "80"
		envHost string
		envPort string
		ok      bool
	)

	if envHost, ok = os.LookupEnv("APP_HOST"); ok {
		host = envHost
	}
	if envPort, ok = os.LookupEnv("APP_PORT"); ok {
		port = envPort
	}

	return host + ":" + port
}

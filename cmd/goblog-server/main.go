package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"

	internalDatabase "github.com/misterabdul/goblog-server/internal/database"
	internalServer "github.com/misterabdul/goblog-server/internal/server"
)

// Run the main server.
func main() {
	var (
		ctx    = context.TODO()
		dbConn *mongo.Database
		server *gin.Engine
		err    error
	)

	if err = godotenv.Load(".env"); err != nil {
		log.Println("Error loading .env file")
	}
	if dbConn, err = internalDatabase.GetDBConnDefault(ctx); err != nil {
		log.Fatal(err)
	}
	server = internalServer.GetServer(dbConn)
	if err = server.Run(getAddress()); err != nil {
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

package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/internal/http/server"
)

// @title                      GoBlog Server
// @version                    1.0
// @description                Simple blog server built with go.
// @BasePath                   /api
// @contact.name               Maintainer
// @contact.email              abdoelrachmad@gmail.com
// @securitydefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                `Bearer <token>`
func main() {
	var (
		ctx            = context.TODO()
		address        string
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
	address = server.ReadContainerAddressFromEnv()
	ginEngine = server.GetServer()
	server.InitRoutes(ginEngine, dbConn, maxCtxDuration)
	server.InitSwagger(ginEngine)
	if err = ginEngine.Run(address); err != nil {
		log.Panic(err)
	}
}

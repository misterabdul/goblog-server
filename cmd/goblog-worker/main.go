package main

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/internal/queue/client"
	"github.com/misterabdul/goblog-server/internal/queue/server"
)

func main() {
	var (
		ctx           = context.TODO()
		dbConn        *mongo.Database
		queueClient   *client.QueueClient
		asynqServer   *asynq.Server
		asynqServeMux *asynq.ServeMux
		err           error
	)

	if err = godotenv.Load(".env"); err != nil {
		log.Println("Error loading .env file")
	}
	if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
		log.Fatal(err)
	}
	queueClient = client.GetClient()
	asynqServer = server.GetServer()
	asynqServeMux = server.InitServeMux(dbConn, queueClient)
	if err = asynqServer.Run(asynqServeMux); err != nil {
		log.Panic(err)
	}
}

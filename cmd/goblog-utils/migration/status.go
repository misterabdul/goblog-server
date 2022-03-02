package migration

import (
	"context"
	"log"

	"github.com/misterabdul/goblog-server/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
)

func Status(ctx context.Context) {
	var (
		dbConn *mongo.Database
		err    error
	)

	if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
		log.Fatal(err)
	}
	defer dbConn.Client().Disconnect(ctx)

	if err = database.Status(ctx, dbConn); err != nil {
		log.Fatal(err)
	}
}

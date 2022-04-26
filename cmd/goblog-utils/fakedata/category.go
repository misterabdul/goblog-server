package fakedata

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
	customMongo "github.com/misterabdul/goblog-server/pkg/mongo"
	"github.com/misterabdul/goblog-server/pkg/utils"
)

func GenerateCategories(ctx context.Context) {
	var (
		dbConn     *mongo.Database
		category   *models.CategoryModel
		categoryId primitive.ObjectID
		now        = primitive.NewDateTimeFromTime(time.Now())
		err        error
	)

	if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
		log.Fatal(err)
	}
	defer dbConn.Client().Disconnect(ctx)
	for i := 0; i < 200; i++ {
		categoryId = primitive.NewObjectID()
		category = &models.CategoryModel{
			UID:       categoryId,
			Slug:      "dummy-cateogry" + fmt.Sprintf("%d", i),
			Name:      "Dummy Category " + fmt.Sprintf("%d", i),
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: nil,
		}
		if err = customMongo.Transaction(ctx, dbConn, false,
			func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
				if sErr = repositories.SaveCategory(
					sCtx,
					dbConn,
					category,
				); sErr != nil {
					return sErr
				}

				return nil
			},
		); err != nil {
			log.Fatal(err)
		}
	}
	utils.ConsolePrintlnGreen("Generated 200 dummy categories.")
}

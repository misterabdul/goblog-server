package fakedata

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/database/repositories"
	customMongo "github.com/misterabdul/goblog-server/pkg/mongo"
	"github.com/misterabdul/goblog-server/pkg/utils"
)

func GenerateCategories(ctx context.Context) {
	var (
		dbConn     *mongo.Database
		repository *repositories.CategoryRepository
		category   *models.CategoryModel
		categoryId primitive.ObjectID
		now        = primitive.NewDateTimeFromTime(time.Now())
		err        error
	)

	if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
		log.Fatal(err)
	}
	defer dbConn.Client().Disconnect(ctx)
	repository = repositories.NewCategoryRepository(dbConn)
	for i := 0; i < 200; i++ {
		categoryId = primitive.NewObjectID()
		category = &models.CategoryModel{
			UID:       categoryId,
			Slug:      "dummy-category" + fmt.Sprintf("%d", i),
			Name:      "Dummy Category " + fmt.Sprintf("%d", i),
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: nil,
		}
		if err = customMongo.Transaction(ctx, dbConn, false,
			func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
				if sErr = repository.Save(
					sCtx, category,
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

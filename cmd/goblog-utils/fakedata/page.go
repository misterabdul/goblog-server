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

func GeneratePages(ctx context.Context) {
	var (
		dbConn            *mongo.Database
		repository        *repositories.PageRepository
		contentRepository *repositories.PageContentRepository
		page              *models.PageModel
		pageContent       *models.PageContentModel
		pageId            primitive.ObjectID
		now               = primitive.NewDateTimeFromTime(time.Now())
		err               error
	)

	if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
		log.Fatal(err)
	}
	defer dbConn.Client().Disconnect(ctx)
	repository = repositories.NewPageRepository(dbConn)
	contentRepository = repositories.NewPageContentRepository(dbConn)
	for i := 0; i < 200; i++ {
		pageId = primitive.NewObjectID()
		page = &models.PageModel{
			UID:         pageId,
			Slug:        fmt.Sprintf("lorem-ipsum-%d", i),
			Title:       fmt.Sprintf("Lorem Ipsum %d", i),
			PublishedAt: randNilOrValue(now),
			CreatedAt:   now,
			UpdatedAt:   now,
			DeletedAt:   nil,
			Author: models.UserCommonModel{
				FirstName: "Super Admin",
				Username:  "superadmin",
				Email:     "superadmin@example.com"}}
		pageContent = &models.PageContentModel{
			UID:     pageId,
			Content: lipsumMarkdown()}
		if err = customMongo.Transaction(ctx, dbConn, false,
			func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
				if sErr = repository.Save(
					sCtx, page,
				); sErr != nil {
					return sErr
				}
				if sErr = contentRepository.Save(
					sCtx, pageContent,
				); sErr != nil {
					return sErr
				}

				return nil
			},
		); err != nil {
			log.Fatal(err)
		}
	}
	utils.ConsolePrintlnGreen("Generated 200 dummy pages.")
}

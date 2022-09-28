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

func GeneratePages(ctx context.Context) {
	var (
		dbConn      *mongo.Database
		page        *models.PageModel
		pageContent *models.PageContentModel
		pageId      primitive.ObjectID
		now         = primitive.NewDateTimeFromTime(time.Now())
		err         error
	)

	if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
		log.Fatal(err)
	}
	defer dbConn.Client().Disconnect(ctx)
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
				if sErr = repositories.SaveOnePage(dbConn, sCtx, page); sErr != nil {
					return sErr
				}
				if sErr = repositories.SaveOnePageContent(dbConn, sCtx, pageContent); sErr != nil {
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

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

func GenerateUsers(ctx context.Context) {
	var (
		dbConn     *mongo.Database
		repository *repositories.UserRepository
		user       *models.UserModel
		userId     primitive.ObjectID
		now        = primitive.NewDateTimeFromTime(time.Now())
		err        error
	)

	if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
		log.Fatal(err)
	}
	defer dbConn.Client().Disconnect(ctx)
	repository = repositories.NewUserRepository(dbConn)
	for i := 0; i < 200; i++ {
		userId = primitive.NewObjectID()
		user = &models.UserModel{
			UID:       userId,
			Username:  "user" + fmt.Sprintf("%d", i),
			Email:     "user" + fmt.Sprintf("%d", i) + "@example.com",
			FirstName: "User " + fmt.Sprintf("%d", i),
			LastName:  "Dummy",
			Password:  "$2y$10$dGPEqemPD8Ys9XQYq3Q.7uVOjeC14tuKYX1iFCzVqqrkJ5B5wX2Xa", // password
			Roles: []models.UserRole{{
				Level: 3,
				Name:  "Writer",
				Since: now}},
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: nil,
		}
		if err = customMongo.Transaction(ctx, dbConn, false,
			func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
				if sErr = repository.Save(
					sCtx, user,
				); sErr != nil {
					return sErr
				}

				return nil
			},
		); err != nil {
			log.Fatal(err)
		}
	}
	utils.ConsolePrintlnGreen("Generated 200 dummy users.")
}

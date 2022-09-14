package me

import (
	"context"

	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/queue/payloads"
	"github.com/misterabdul/goblog-server/internal/service"
)

func UpdateMe(dbConn *mongo.Database) (handler asynq.HandlerFunc) {
	return func(ctx context.Context, t *asynq.Task) error {
		var (
			postService = service.NewPostService(nil, ctx, dbConn)
			payload     *payloads.UpdateMePayload
			err         error
		)

		if payload, err = payloads.UnmarshallUpdateMePayload(t.Payload()); err != nil {
			return err
		}
		if err = postService.UpdatePostAuthor(&payload.UserModel); err != nil {
			return err
		}

		return nil
	}
}

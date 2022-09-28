package me

import (
	"context"

	"github.com/hibiken/asynq"

	"github.com/misterabdul/goblog-server/internal/queue/payloads"
	"github.com/misterabdul/goblog-server/internal/service"
)

func UpdateMe(
	svc *service.Service,
) (handler asynq.HandlerFunc) {

	return func(ctx context.Context, t *asynq.Task) error {
		var (
			payload *payloads.UpdateMePayload
			err     error
		)

		if payload, err = payloads.UnmarshallUpdateMePayload(t.Payload()); err != nil {
			return err
		}
		if err = svc.Post.UpdateManyAuthor(ctx, &payload.UserModel); err != nil {
			return err
		}

		return nil
	}
}

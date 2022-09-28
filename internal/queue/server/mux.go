package server

import (
	"context"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/queue"
	"github.com/misterabdul/goblog-server/internal/queue/client"
	meHandler "github.com/misterabdul/goblog-server/internal/queue/handlers/me"
	"github.com/misterabdul/goblog-server/internal/service"
)

func InitServeMux(
	dbConn *mongo.Database,
	queueClient *client.QueueClient,
) *asynq.ServeMux {
	var (
		mux              *asynq.ServeMux
		serverRelatedEnv = getRedisServerRelatedEnv()
		svc              = service.NewService(dbConn, queueClient)
	)

	mux = asynq.NewServeMux()
	switch serverRelatedEnv.Mode {
	default:
		fallthrough
	case 0:
		break
	case 1:
		break
	case 2:
		mux.Use(loggingMiddleware)
	}

	mux.HandleFunc(queue.UpdateMe, meHandler.UpdateMe(svc))

	return mux
}

func loggingMiddleware(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(
		ctx context.Context,
		t *asynq.Task,
	) (err error) {
		var start = time.Now()

		log.Printf("Start processing \"%s\"\n", t.Type())
		if err = h.ProcessTask(ctx, t); err != nil {
			return err
		}
		log.Printf("Finished processing \"%s\": Elapsed Time = %v", t.Type(), time.Since(start))

		return nil
	})
}

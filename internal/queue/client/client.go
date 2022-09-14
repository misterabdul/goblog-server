package client

import (
	"encoding/json"

	"github.com/hibiken/asynq"

	"github.com/misterabdul/goblog-server/internal/queue/server"
)

type QueueClient struct {
	opts   asynq.RedisClientOpt
	client *asynq.Client
}

func (c *QueueClient) Connect() {
	c.client = asynq.NewClient(c.opts)
}

func (c *QueueClient) Disconnect() (err error) {
	return c.client.Close()
}

func (c *QueueClient) NewTask(
	taskName string,
	payload interface{},
	opts ...asynq.Option,
) (err error) {
	var (
		task     *asynq.Task
		_payload []byte
	)

	if _payload, err = json.Marshal(payload); err != nil {
		return err
	}
	task = asynq.NewTask(taskName, _payload, opts...)
	if _, err = c.client.Enqueue(task); err != nil {
		return err
	}

	return nil
}

func GetClient() (client *QueueClient) {
	return &QueueClient{
		opts: server.ReadRedisOptsFromEnv()}
}

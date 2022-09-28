package client

import (
	"encoding/json"
	"os"

	"github.com/hibiken/asynq"
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
	client = &QueueClient{
		opts: readRedisOptsFromEnv()}

	client.Connect()

	return client
}

func readRedisOptsFromEnv() asynq.RedisClientOpt {
	var (
		host    = "localhost"
		port    = "6379"
		user    = ""
		pass    = ""
		envHost string
		envPort string
		envUser string
		envPass string
		ok      bool
	)

	if envHost, ok = os.LookupEnv("REDIS_HOST"); ok {
		host = envHost
	} else {
		host = "localhost"
	}
	if envPort, ok = os.LookupEnv("REDIS_PORT"); ok {
		port = envPort
	} else {
		port = "80"
	}
	if envUser, ok = os.LookupEnv("REDIS_USER"); ok {
		user = envUser
	} else {
		user = ""
	}
	if envPass, ok = os.LookupEnv("REDIS_PASS"); ok {
		pass = envPass
	} else {
		pass = "80"
	}

	return asynq.RedisClientOpt{
		Addr:     host + ":" + port,
		Username: user,
		Password: pass,
	}
}

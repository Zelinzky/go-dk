package integrationtest

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/maragudk/env"

	"go-dk/messaging"
)

// CreateQueue for testing.
// Usage:
//
//	queue, cleanup := CreateQueue()
//	defer cleanup()
//	...
func CreateQueue() (*messaging.Queue, func()) {
	env.MustLoad("../.env-test")

	queueName := env.GetStringOrDefault("QUEUE_NAME", "jons")
	queue := messaging.NewQueue(messaging.NewQueueOptions{
		Config: getAWSConfig(),
		Name:   queueName,
	})

	createQueueOutput, err := queue.Client.CreateQueue(context.Background(), &sqs.CreateQueueInput{
		QueueName: &queueName,
	})
	if err != nil {
		panic(err)
	}

	return queue, func() {
		_, err = queue.Client.DeleteQueue(context.Background(), &sqs.DeleteQueueInput{
			QueueUrl: createQueueOutput.QueueUrl,
		})
		if err != nil {
			panic(err)
		}
	}
}

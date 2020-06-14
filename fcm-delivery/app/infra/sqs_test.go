package infra

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSQSRepository_DeleteMessageBatch(t *testing.T) {
	sqsClient, sqsEvent := setup(t)
	defer teardown(t, sqsClient)

	sqsRepository := NewSQSRepository(sqsClient, *queueURL)

	deleteMessageBatchOutput, err := sqsRepository.DeleteMessageBatch(sqsEvent)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, deleteMessageBatchOutput.Failed, 0)
	assert.Len(t, deleteMessageBatchOutput.Successful, 1)
}

func setup(t *testing.T) (*sqs.SQS, events.SQSEvent) {
	// TODO run "localstack start" in advance.

	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String(sqsEndpoint),
		Region:   aws.String("us-east-1"),
	}))
	sqsClient := sqs.New(sess)

	createQueueOutput, err := setupCreateQueue(t, sqsClient)
	if err != nil {
		t.Fatal(err)
	}

	queueURL = createQueueOutput.QueueUrl
	_, err = setupSendMessage(t, sqsClient)
	if err != nil {
		t.Fatal(err)
	}

	sqsEvent, err := setupCreateSQSEvent(t, sqsClient)
	if err != nil {
		t.Fatal(err)
	}

	return sqsClient, *sqsEvent
}

func setupCreateQueue(t *testing.T, sqsClient *sqs.SQS) (*sqs.CreateQueueOutput, error) {
	return sqsClient.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	})
}

func setupSendMessage(t *testing.T, sqsClient *sqs.SQS) (*sqs.SendMessageBatchOutput, error) {
	return sqsClient.SendMessageBatch(&sqs.SendMessageBatchInput{
		Entries: []*sqs.SendMessageBatchRequestEntry{
			{
				Id:          aws.String("1"),
				MessageBody: aws.String("msg_1"),
			},
		},
		QueueUrl: queueURL,
	})
}

func setupCreateSQSEvent(t *testing.T, sqsClient *sqs.SQS) (*events.SQSEvent, error) {
	receiveMessageOutput, err := receiveMessage(sqsClient)
	if err != nil {
		return nil, err
	}

	var records []events.SQSMessage
	for _, message := range receiveMessageOutput.Messages {
		sqsMessage := events.SQSMessage{
			MessageId:     *message.MessageId,
			ReceiptHandle: *message.ReceiptHandle,
			Body:          *message.Body,
		}
		records = append(records, sqsMessage)
	}

	return &events.SQSEvent{Records: records}, nil
}

func receiveMessage(sqsClient *sqs.SQS) (*sqs.ReceiveMessageOutput, error) {
	return sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl: queueURL,
	})
}

func teardown(t *testing.T, sqsClient *sqs.SQS) {
	t.Log("🌴")
	_, err := sqsClient.DeleteQueue(&sqs.DeleteQueueInput{QueueUrl: queueURL})

	if err != nil {
		t.Fatal(err)
	}
}

const sqsEndpoint = "http://localhost:4576"
const queueName = "testqueue"

var queueURL *string

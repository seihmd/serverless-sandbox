package infra

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSRepository struct {
	client   *sqs.SQS
	queueURL string
}

func NewSQSRepository(client *sqs.SQS, queueURL string) *SQSRepository {
	return &SQSRepository{client: client, queueURL: queueURL}
}

func (r *SQSRepository) DeleteMessageBatch(e events.SQSEvent) (*sqs.DeleteMessageBatchOutput, error) {
	var entries []*sqs.DeleteMessageBatchRequestEntry
	for _, r := range e.Records {
		entries = append(entries, &sqs.DeleteMessageBatchRequestEntry{
			Id:            aws.String(r.MessageId),
			ReceiptHandle: aws.String(r.ReceiptHandle),
		})
	}

	return r.client.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
		Entries:  entries,
		QueueUrl: aws.String(r.queueURL),
	})
}

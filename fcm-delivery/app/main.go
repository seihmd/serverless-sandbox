package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/serverless-sandbox/fcm-delivery/app/application"
	"github.com/serverless-sandbox/fcm-delivery/app/domain"
	"github.com/serverless-sandbox/fcm-delivery/app/infra"
	"os"
)

type Response events.APIGatewayProxyResponse

var firebaseServiceCredJson string
var sqsClient *sqs.SQS
var queueURL string

func init() {
	firebaseServiceCredJson = os.Getenv("FIREBASE_CREDENTIAL")

	sess := session.Must(session.NewSession())
	sqsClient = sqs.New(sess)

	queueURL = os.Getenv("SQS_QUEUE_URL")
}

func Handler(ctx context.Context, e events.SQSEvent) (Response, error) {
	messages := application.CreateMessages(e)

	m, err := infra.NewFCM(ctx, firebaseServiceCredJson)
	if err != nil {
		return Response{
			StatusCode: 500,
			Body:       err.Error(),
		}, err
	}

	rr := infra.ResultRecorder{}
	p := domain.NewPushDelivery(m, rr)

	err = p.Handle(ctx, messages)
	if err != nil {
		return Response{
			StatusCode: 500,
			Body:       err.Error(),
		}, err
	}

	sqsRepository := infra.NewSQSRepository(sqsClient, queueURL)
	output, err := sqsRepository.DeleteMessageBatch(e)

	if err != nil {
		return Response{
			StatusCode: 500,
			Body:       err.Error(),
		}, err
	}

	return Response{
		Body: output.String(),
	}, err
}

func main() {
	lambda.Start(Handler)
}

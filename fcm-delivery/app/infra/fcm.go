package infra

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/serverless-sandbox/fcm-delivery/app/domain"
	"google.golang.org/api/option"
)

type FCM struct {
	client *messaging.Client
}

func (fcm FCM) Send(ctx context.Context, messages []domain.Message) (*messaging.BatchResponse, error) {
	var fcmMessages []*messaging.Message

	for _, m := range messages {
		if fcmMessage := buildMessage(m); fcmMessage != nil {
			fcmMessages = append(fcmMessages, fcmMessage)
		}
	}

	return fcm.client.SendAll(ctx, fcmMessages)
}

func NewFCM(ctx context.Context, credentialsJson string) (*FCM, error) {
	clientOption := option.WithCredentialsJSON([]byte(credentialsJson))

	app, err := firebase.NewApp(ctx, nil, clientOption)
	if err != nil {
		return nil, err
	}

	fcmClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	return &FCM{
		client: fcmClient,
	}, nil
}

package main

import (
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"google.golang.org/api/option"
	"log"
	"os"
)

type Platform int

var (
	iOS     Platform = 1
	android Platform = 2
	web     Platform = 3
)

var fcmClient *messaging.Client

type deliveryItem struct {
	Title    string   `json:"title"`
	Body     string   `json:"body"`
	Platform Platform `json:"platform"`
	Token    string   `json:"token"`
}

type Response events.APIGatewayProxyResponse

func init() {
	cred := os.Getenv("FIREBASE_CREDENTIAL")

	log.Println(cred)

	sa := option.WithCredentialsJSON([]byte(cred))
	log.Println(sa)

	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fcmClient, err = app.Messaging(ctx)
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func Handler(ctx context.Context, e events.SQSEvent) (Response, error) {
	log.Println("🌴")
	var messages []*messaging.Message
	for _, message := range e.Records {
		log.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)

		if m := createMessage(message.Body); m != nil {
			messages = append(messages, m)
		}
	}

	batchResponse, err := fcmClient.SendAll(ctx, messages)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(batchResponse)
	if err != nil {
		log.Fatal(err)
	}

	res := Response{
		Body: string(b),
	}

	return res, nil
}

func createMessage(body string) *messaging.Message {
	var item deliveryItem
	if err := json.Unmarshal([]byte(body), &item); err != nil {
		log.Fatal(err)
	}

	if item.Platform == iOS || item.Platform == android {
		// TODO implement
		return nil
	}

	n := messaging.Notification{
		Title: item.Title,
		Body:  item.Body,
	}

	m := messaging.Message{
		//Data:         nil,
		Notification: &n,
		//Android:      nil,
		//Webpush:      nil,
		//APNS:         nil,
		//FCMOptions:   nil,
		Token: item.Token,
		//Topic:        "",
		//Condition:    "",
	}

	return &m
}

func main() {
	lambda.Start(Handler)
}

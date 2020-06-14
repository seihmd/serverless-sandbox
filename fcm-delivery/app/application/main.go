package application

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/serverless-sandbox/fcm-delivery/app/domain"
	"log"
)

func CreateMessages(e events.SQSEvent) []*domain.Message {
	var messages []*domain.Message
	for _, record := range e.Records {
		log.Printf("The record %s for event source %s = %s \n", record.MessageId, record.EventSource, record.Body)

		var m domain.Message
		if err := json.Unmarshal([]byte(record.Body), &m); err != nil {
			log.Println(err)
		} else {
			messages = append(messages, &m)
		}
	}

	return messages
}

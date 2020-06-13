package infra

import (
	"firebase.google.com/go/messaging"
	"log"
)

type ResultRecorder struct {
}

func (r ResultRecorder) ErrorOccurred(err error) {
	log.Printf("error occurred: %s", err.Error())
}

func (r ResultRecorder) Completed(b *messaging.BatchResponse) {
	log.Println("completed sending.")
	log.Printf("failure count: %d", b.FailureCount)
	log.Printf("success count: %d", b.SuccessCount)

	if b.FailureCount == 0 {
		return
	}

	for _, r := range b.Responses {
		if r.Error != nil {
			log.Printf("messageId: %s failed. error: %s\n", r.MessageID, r.Error.Error())
		}
	}
}

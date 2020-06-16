package domain

import (
	"context"
	"firebase.google.com/go/messaging"
)

type Platform int

var (
	IOS     Platform = 1
	Android Platform = 2
	Web     Platform = 3
)

type Message struct {
	Title          string
	Body           string
	ImageURL       string `json:"image_url"`
	Token          string
	Platform       Platform
	AnalyticsLabel string `json:"analytics_label"`
}

func (m Message) IsIOS() bool {
	return m.Platform == IOS
}

func (m Message) IsAndroid() bool {
	return m.Platform == Android
}

func (m Message) IsWeb() bool {
	return m.Platform == Web
}

type MessageSender interface {
	Send(context.Context, []*Message) (*messaging.BatchResponse, error)
}

type ResultRecorder interface {
	ErrorOccurred(error)
	Completed(*messaging.BatchResponse)
}

type PushDelivery struct {
	messageSender  MessageSender
	resultRecorder ResultRecorder
}

func NewPushDelivery(ms MessageSender, rr ResultRecorder) *PushDelivery {
	return &PushDelivery{
		messageSender:  ms,
		resultRecorder: rr,
	}
}

func (s *PushDelivery) Handle(ctx context.Context, m []*Message) error {
	batchResponse, err := s.messageSender.Send(ctx, m)

	s.logResult(batchResponse, err)

	return err
}

func (s *PushDelivery) logResult(b *messaging.BatchResponse, err error) {
	if err != nil {
		s.resultRecorder.ErrorOccurred(err)
		return
	}

	s.resultRecorder.Completed(b)
}

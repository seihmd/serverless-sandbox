package application

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateMessages_validMessageBody_createMessageFromIt(t *testing.T) {
	e := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId: "1",
				Body:      `{}`,
			},
			{
				MessageId: "2",
				Body:      `{}`,
			},
			{
				MessageId: "3",
				Body:      `{}`,
			},
		},
	}

	messages := CreateMessages(e)
	assert.Len(t, messages, 3)
}

func TestCreateMessages_invalidMessageBody_ignoreIt(t *testing.T) {
	e := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId: "1",
				Body:      `{}`,
			},
			{
				MessageId: "2",
				Body:      `INVALID`,
			},
			{
				MessageId: "3",
				Body:      `{}`,
			},
		},
	}

	messages := CreateMessages(e)
	assert.Len(t, messages, 2)
}

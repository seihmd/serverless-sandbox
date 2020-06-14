package domain

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestPushDeliveryService_Handle_errorOccurred_recordErrorOtherwiseRecordCompleted(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		expect int
	}{
		{
			name:   "no error",
			err:    nil,
			expect: 0,
		},
		{
			name:   "error occurred",
			err:    errors.New("test"),
			expect: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			messageSender := NewMockMessageSender(ctrl)
			resultRecorder := NewMockResultRecorder(ctrl)

			s := PushDelivery{
				messageSender:  messageSender,
				resultRecorder: resultRecorder,
			}

			messageSender.
				EXPECT().
				send(gomock.Any(), gomock.Any()).
				Return(nil, tt.err)

			resultRecorder.
				EXPECT().
				errorOccurred(tt.err).
				Times(tt.expect)

			resultRecorder.
				EXPECT().
				completed(tt.err).
				Times(1 - tt.expect)

			s.Handle(context.Background(), []Message{})
		})
	}
}

func TestMessage_jsonString_canUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		jsonMessage string
		expect      Message
	}{
		{
			name: "web push message",
			jsonMessage: `{"title": "タイトル",` +
				`"body": "本文",` +
				`"image_url": "https://example.com",` +
				`"token":"token",` +
				`"platform": 3,` +
				`"analytics_label": "analytics_label"}`,
			expect: Message{
				Title:          "タイトル",
				Body:           "本文",
				ImageURL:       "https://example.com",
				Token:          "token",
				Platform:       Web,
				AnalyticsLabel: "analytics_label",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var m Message
			err := json.Unmarshal([]byte(tt.jsonMessage), &m)

			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expect, m)
		})
	}
}

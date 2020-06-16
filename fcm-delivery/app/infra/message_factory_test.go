package infra

import (
	"firebase.google.com/go/messaging"
	"github.com/serverless-sandbox/fcm-delivery/app/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_buildMessage(t *testing.T) {
	type args struct {
		m *domain.Message
	}
	tests := []struct {
		name string
		args args
		want *messaging.Message
	}{
		{
			"web push message",
			args{
				&domain.Message{
					Title:          "タイトル1",
					Body:           "本文1",
					ImageURL:       "https://example.com",
					Token:          "test_token_1",
					Platform:       domain.Web,
					AnalyticsLabel: "analytics_label_1",
				},
			},
			&messaging.Message{
				Webpush: &messaging.WebpushConfig{
					Notification: &messaging.WebpushNotification{
						Title: "タイトル1",
						Body:  "本文1",
						Image: "https://example.com",
					},
				},
				FCMOptions: &messaging.FCMOptions{
					AnalyticsLabel: "analytics_label_1",
				},
				Token: "test_token_1",
			},
		},
		{
			"ios push message",
			args{
				&domain.Message{
					Title:          "タイトル2",
					Body:           "本文2",
					ImageURL:       "https://example.com",
					Token:          "test_token_2",
					Platform:       domain.IOS,
					AnalyticsLabel: "analytics_label_2",
				},
			},
			&messaging.Message{
				APNS: &messaging.APNSConfig{
					Payload: &messaging.APNSPayload{
						Aps: &messaging.Aps{
							Alert: &messaging.ApsAlert{
								Title: "タイトル2",
								Body:  "本文2",
							},
						},
					},
					FCMOptions: &messaging.APNSFCMOptions{
						AnalyticsLabel: "analytics_label_2",
						ImageURL:       "https://example.com",
					},
				},
				Token: "test_token_2",
			},
		},
		{
			"android push message",
			args{
				&domain.Message{
					Title:          "タイトル3",
					Body:           "本文3",
					ImageURL:       "https://example.com",
					Token:          "test_token_3",
					Platform:       domain.Android,
					AnalyticsLabel: "analytics_label_3",
				},
			},
			&messaging.Message{
				Android: &messaging.AndroidConfig{
					Notification: &messaging.AndroidNotification{
						Title:    "タイトル3",
						Body:     "本文3",
						ImageURL: "https://example.com",
					},
					FCMOptions: &messaging.AndroidFCMOptions{
						AnalyticsLabel: "analytics_label_3",
					},
				},
				Token: "test_token_3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := buildMessage(tt.args.m)
			assert.Equal(t, tt.want, m)
		})
	}
}

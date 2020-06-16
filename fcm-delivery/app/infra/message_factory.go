package infra

import (
	"firebase.google.com/go/messaging"
	"github.com/serverless-sandbox/fcm-delivery/app/domain"
)

// see https://firebase.google.com/docs/reference/fcm/rest/v1/projects.messages

func buildMessage(m *domain.Message) *messaging.Message {
	if m.Platform == domain.Web {
		return buildWebPushMessage(m)
	} else if m.Platform == domain.IOS {
		return buildIOSPushMessage(m)
	} else if m.Platform == domain.Android {
		return buildAndroidPushMessage(m)
	}

	return nil
}

func buildWebPushMessage(m *domain.Message) *messaging.Message {
	if m.Platform != domain.Web {
		return nil
	}

	return &messaging.Message{
		Webpush: &messaging.WebpushConfig{
			Notification: &messaging.WebpushNotification{
				Title: m.Title,
				Body:  m.Body,
				Image: m.ImageURL,
			},
		},
		FCMOptions: &messaging.FCMOptions{
			AnalyticsLabel: m.AnalyticsLabel,
		},
		Token: m.Token,
	}
}

func buildIOSPushMessage(m *domain.Message) *messaging.Message {
	if m.Platform != domain.IOS {
		return nil
	}

	return &messaging.Message{
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: m.Title,
						Body:  m.Body,
					},
				},
			},
			FCMOptions: &messaging.APNSFCMOptions{
				AnalyticsLabel: m.AnalyticsLabel,
				ImageURL:       m.ImageURL,
			},
		},
		Token: m.Token,
	}
}

func buildAndroidPushMessage(m *domain.Message) *messaging.Message {
	if m.Platform != domain.Android {
		return nil
	}

	return &messaging.Message{
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				Title:    m.Title,
				Body:     m.Body,
				ImageURL: m.ImageURL,
			},
			FCMOptions: &messaging.AndroidFCMOptions{
				AnalyticsLabel: m.AnalyticsLabel,
			},
		},
		Token: m.Token,
	}
}

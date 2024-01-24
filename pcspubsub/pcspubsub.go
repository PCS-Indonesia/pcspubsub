package pcspubsub

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

type CommandMessage struct {
	Command string `json:"command"`
	Payload string `json:"payload"`
	ID      uint   `json:"id"`
	Detail  string `json:"detail"`
}

// Set your Google Cloud project ID and topic name
// projectID := "pcs-drive-350809"
// topicName := "test-topic"
// subscriptionName := "user-data-subscription"
// credentialsPath := "cred.json"

func DispatchSub(projectID string, topicName string, subscriptionName string, credentialsPath string, callback func(msg CommandMessage)) error {

	// Create a new Pub/Sub client with service account credentials
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return err
	}

	// Create a new subscription
	subscription := client.Subscription(subscriptionName)
	exists, err := subscription.Exists(ctx)
	if err != nil {
		return err
	}
	if !exists {
		topic := client.Topic(topicName)
		subscription, err = client.CreateSubscription(ctx, subscriptionName, pubsub.SubscriptionConfig{
			Topic: topic,
		})
		if err != nil {
			return err
		}
	}

	// Start receiving messages
	err = subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		var cmd CommandMessage
		rs := json.Unmarshal(msg.Data, &cmd)
		if rs == nil {
			callback(cmd)
			msg.Ack()
		}
	})

	if err != nil {
		return err
	}

	return nil
}

func Publish(projectID string, topicName string, subscriptionName string, credentialsPath string, msg CommandMessage) error {
	// Create a new context and client
	ctx := context.Background()
	// Create a new Pub/Sub client with service account credentials
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return err
	}

	topic := client.Topic(topicName)
	if ok, _ := topic.Exists(ctx); !ok {
		topic, err = client.CreateTopic(ctx, topicName)
		if err != nil {
			return err
		}
	}

	// Publish a text message on the created topic
	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	topic.Publish(ctx, &pubsub.Message{
		Data: message,
	})

	return nil
}

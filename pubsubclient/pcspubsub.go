package pubsubclient

import (
	"context"
	"encoding/json"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

type CommandMessage struct {
	Command string          `json:"command"`
	Payload string          `json:"payload"`
	ID      uint            `json:"id"`
	Detail  string          `json:"detail"`
	Message *pubsub.Message `json:"-"`
}

// Set your Google Cloud project ID and topic name
// projectID := "pcs-drive-350809"
// topicName := "test-topic"
// subscriptionName := "user-data-subscription"
// credentialsPath := "cred.json"

type PubSubClient struct {
	ctx                   context.Context
	client                *pubsub.Client
	maxConcurrentMessages int
}

func NewPubSubClient(ctx context.Context, projectID string, credentialsPath string, maxConcurrent int) (*PubSubClient, error) {
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return nil, err
	}
	if maxConcurrent == 0 {
		maxConcurrent = 1
	}

	return &PubSubClient{
		ctx:                   ctx,
		client:                client,
		maxConcurrentMessages: maxConcurrent,
	}, nil
}

func (c *PubSubClient) ReceiveMessages(subscriptionName string, callback func(ctx context.Context, msg CommandMessage) error) error {

	subscription := c.client.Subscription(subscriptionName)
	subscription.ReceiveSettings.MaxOutstandingMessages = c.maxConcurrentMessages
	err := subscription.Receive(c.ctx, func(ctx context.Context, msg *pubsub.Message) {
		timeout := 120 * time.Second
		ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
		defer cancel() // Ensure the context is cancelled when the function returns
		var cmd CommandMessage
		err := json.Unmarshal(msg.Data, &cmd)
		if err == nil {
			cmd.Message = msg
			if err := callback(ctxWithTimeout, cmd); err == nil {
				msg.Ack()
			} else {
				msg.Nack()
				return
			}
		} else {
			msg.Nack()
			return
		}
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *PubSubClient) PublishMessage(topicName string, msg CommandMessage) error {
	topic := c.client.Topic(topicName)
	if ok, _ := topic.Exists(c.ctx); !ok {
		_, err := c.client.CreateTopic(c.ctx, topicName)
		if err != nil {
			return err
		}
	}

	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = topic.Publish(c.ctx, &pubsub.Message{
		Data: message,
	}).Get(c.ctx)

	return err
}

func (c *PubSubClient) CreateSubscription(subscriptionName string, topicName string) (*pubsub.Subscription, error) {
	return c.client.CreateSubscription(c.ctx, subscriptionName, pubsub.SubscriptionConfig{
		Topic: c.client.Topic(topicName),
	})
}

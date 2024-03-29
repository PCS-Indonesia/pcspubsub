package pubsubclient

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

type PubSubClient struct {
	ctx    context.Context
	client *pubsub.Client
}

func NewPubSubClient(projectID string, credentialsPath string) (*PubSubClient, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return nil, err
	}

	return &PubSubClient{
		ctx:    ctx,
		client: client,
	}, nil
}

func (c *PubSubClient) ReceiveMessages(subscriptionName string, callback func(msg CommandMessage)) error {
	subscription := c.client.Subscription(subscriptionName)
	err := subscription.Receive(c.ctx, func(ctx context.Context, msg *pubsub.Message) {
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

	topic.Publish(c.ctx, &pubsub.Message{
		Data: message,
	})

	return nil
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/PCS-Indonesia/pcspubsub/pubsubclient"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("can't load env")
	}
	os.Setenv("ID", "Asia/Jakarta")
}

func main() {
	var arg string
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}
	switch arg {
	case "pub":
		pub()
	case "pubv2":
		pubV2()
	case "script":
		createSub()
	default:
		sub()
	}
}

func pub() {
	var message pubsubclient.CommandMessage
	message.Command = "notify"
	message.Payload = "lorem ipsum dolor sit amet"
	message.ID = 0
	message.Detail = "john doe"

	topic := os.Getenv("PUBSUB_MMS_TOPIC")

	ctx := context.TODO()
	client, err := pubsubclient.NewPubSubClient(ctx, os.Getenv("PUBSUB_PROJECT_ID"), os.Getenv("PUBSUB_CREDENTIAL"), 1)
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}

	err = client.PublishMessage(topic, message)
	fmt.Printf("Message: %+v\n", message)

	if err != nil {
		log.Printf("System: Error publish message topic %s, error: %v", topic, err)
	}
}

func pubV2() {
	var message pubsubclient.CommandMessage
	message.Command = "notify"
	message.Payload = "from external Account"
	message.ID = 0
	message.Detail = "from external Account"

	topic := os.Getenv("PUBSUB_MMS_TOPIC")

	ctx := context.TODO()
	pubsub := pubsubclient.PubSubConfig{
		ProjectID:     os.Getenv("PUBSUB_PROJECT_ID"),
		TokenSource:   os.Getenv("PUBSUB_CREDENTIAL"),
		MaxConcurrent: 1,
		ExpiredToken:  time.Now().Add(1 * time.Hour), // Set expiration to 1 hour from now
		Ctx:           ctx,
	}

	client, err := pubsub.NewPubSubClientWithTokenWIF()
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
		return
	}

	err = client.PublishMessage(topic, message)
	fmt.Printf("Message: %+v\n", message)

	if err != nil {
		log.Printf("System: Error publish message topic %s, error: %v", topic, err)
	}
}

func createSub() {
	ctx := context.TODO()
	client, err := pubsubclient.NewPubSubClient(ctx, os.Getenv("PUBSUB_PROJECT_ID"), os.Getenv("PUBSUB_CREDENTIAL"), 1)
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}

	if _, err := client.CreateSubscription(os.Getenv("PUBSUB_MMS_SUBS"), os.Getenv("PUBSUB_MMS_TOPIC")); err != nil {
		log.Printf("Failed to create subscription: %v", err)
	} else {
		log.Println("Subscription created")
	}
}

func sub() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	client, err := pubsubclient.NewPubSubClient(ctx, os.Getenv("PUBSUB_PROJECT_ID"), os.Getenv("PUBSUB_CREDENTIAL"), 1)
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		sub := os.Getenv("PUBSUB_MMS_SUBS")
		err = client.ReceiveMessages(sub, subFunction)

		if err != nil {
			log.Panicf("Failed to receive messages from : %s, %v\n", sub, err)
		}

		log.Println("System has shutdown gracefully")
	}()

	// Cancel the context when a signal is received
	<-sigchan
	cancel()

	// Wait for the goroutine to finish
	wg.Wait()
}

func subFunction(ctx context.Context, msg pubsubclient.CommandMessage) error {
	fmt.Printf("Message: %+v\n", msg)
	return nil
}

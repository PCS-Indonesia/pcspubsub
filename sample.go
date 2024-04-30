package main

import (
	"fmt"
	"log"
	"os"
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
	default:
		sub()

		for {
			time.Sleep(5000 * time.Millisecond)
		}
	}
}

func pub() {
	var message pubsubclient.CommandMessage
	message.Command = "notify"
	message.Payload = "lorem ipsum dolor sit amet"
	message.ID = 0
	message.Detail = "john doe"

	topic := os.Getenv("PUBSUB_MMS_TOPIC")

	client, err := pubsubclient.NewPubSubClient(os.Getenv("PUBSUB_PROJECT_ID"), os.Getenv("PUBSUB_CREDENTIAL"))
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}

	err = client.PublishMessage(topic, message)
	fmt.Printf("Message: %+v\n", message)

	if err != nil {
		log.Printf("System: Error publish message topic %s, error: %v", topic, err)
	}
}

func sub() {
	client, err := pubsubclient.NewPubSubClient(os.Getenv("PUBSUB_PROJECT_ID"), os.Getenv("PUBSUB_CREDENTIAL"))
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}

	go func() {
		sub := os.Getenv("PUBSUB_MMS_SUBS")
		err = client.ReceiveMessages(sub, subFunction)

		if err != nil {
			log.Printf("Failed to receive messages from : %s, %v", sub, err)
		}
	}()
}

func subFunction(msg pubsubclient.CommandMessage) {
	fmt.Printf("Message: %+v\n", msg)
}

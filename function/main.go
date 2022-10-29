// Package helloworld provides a set of Cloud Functions samples.
package helloworld

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"google.golang.org/api/compute/v1"
)

func init() {
	functions.CloudEvent("HelloPubSub", helloPubSub)
}

// MessagePublishedData contains the full Pub/Sub message
// See the documentation for more details:
// https://cloud.google.com/eventarc/docs/cloudevents#pubsub
type MessagePublishedData struct {
	Message PubSubMessage
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// helloPubSub consumes a CloudEvent message and extracts the Pub/Sub message.
func helloPubSub(ctx context.Context, e event.Event) error {
	cclient, err := compute.NewService(ctx)
	if err != nil {
		return fmt.Errorf("error")
	}
	is := cclient.Instances.List(os.Getenv("PROJECT"),os.Getenv("LOCATION"))
	ts, err := is.Do()
	if err != nil {
		return fmt.Errorf("error")
	}
	log.Printf("%#v", ts.Items)

	if len(ts.Items) >= 3 {
		pclient, err := pubsub.NewClient(ctx,os.Getenv("PROJECT"))
		if err != nil {
			return fmt.Errorf("error")
		}
		defer pclient.Close()
		t := pclient.Topic(os.Getenv("TOPIC_WAITING"))
		r := t.Publish(ctx, &pubsub.Message{
			//Data: []byte("test1"),
			Data: e.Data(),
		})

		id, err := r.Get(ctx)
		if err != nil {
			return fmt.Errorf("error")
		}
		fmt.Println(id)
	}


	// just showing
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %v", err)
	}

	name := string(msg.Message.Data) // Automatically decoded from base64.
	if name == "" {
		name = "World"
	}
	log.Printf("Hello, %s!", name)
	return nil
}

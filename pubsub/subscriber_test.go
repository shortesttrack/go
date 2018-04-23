package pubsub

import (
	"testing"
	"st-go/util"
	"context"
	"log"
)

func TestSubscriber_Start(t *testing.T) {
	projectID, err := util.MustGetEnv("PUBSUB_PROJECT_ID")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	ctx, _ = context.WithCancel(ctx)
	//defer cancel()

	s, err := NewSubscriber(ctx, SubscriptionOptions{
		ProjectID: projectID,
		Subscription: "not-e-sub-tets",
	})
	if err != nil {
		t.Fatal(err)
	}

	msgCh := s.Start()
	for {
		select {
		case <-ctx.Done():
			log.Printf("context done error: %v", ctx.Err())
		case msg, ok := <-msgCh:
			if !ok {
				return
			}
			log.Printf("msg from channel: %v", string(msg.Data()))
			msg.Ack()
			log.Printf("acked")
		}
	}
}
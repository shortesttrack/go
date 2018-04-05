package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"st-go/errors"
)

type PublishOptions struct {
	ProjectID string
	Topic string
}

func NewPublisher(options PublishOptions) (Publisher, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, options.ProjectID)
	if err != nil {
		return nil, err
	}
	topic := client.Topic(options.Topic)
	exist, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New(`invalid topic: ` + options.Topic)
	}
	return &publisher{topic}, nil
}

type publisher struct {
	*pubsub.Topic
}

func (p *publisher) Publish(ctx context.Context, message Message) error {
	result := p.Topic.Publish(ctx, &pubsub.Message{
		Data: message.Data(),
		Attributes: message.Attributes(),
	})
	_, err := result.Get(ctx)
	return err
}
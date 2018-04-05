package pubsub

import (
	"context"
	"cloud.google.com/go/pubsub"
	"st-go/errors"
	"sync"
)

type SubscriptionOptions struct {
	ProjectID    string
	Subscription string
}

func NewSubscriber(ctx context.Context, options SubscriptionOptions) (Subscriber, error) {
	client, err := pubsub.NewClient(ctx, options.ProjectID)
	if err != nil {
		return nil, err
	}

	subscription := client.Subscription(options.Subscription)
	exists, err := subscription.Exists(ctx)
	if err != nil {
		client.Close()
		return nil, err
	}
	if exists == false {
		client.Close()
		return nil, errors.New(`invalid subscription: ` + options.Subscription)
	}
	return &subscriber{sub: subscription, ctx: ctx}, nil
}

type subscriber struct {
	sub    *pubsub.Subscription
	ctx    context.Context
	cancel func()
	err    error

	stopMu sync.Mutex
	stopped bool
}

func (s *subscriber) Start() <-chan Message {
	output := make(chan Message)
	go func(s *subscriber, mc chan Message) {
		defer close(output)
		s.ctx, s.cancel = context.WithCancel(s.ctx)
		err := s.sub.Receive(s.ctx, func(ctx context.Context, message *pubsub.Message) {
			output <- &message{message}
		})
		if err != nil {
			s.Stop()
			s.err = err
		}
	}(s, output)
	return output
}

func (s *subscriber) Err() error {
	return s.err
}

func (s *subscriber) Stop() error {
	s.stopMu.Lock()
	defer s.stopMu.Unlock()
	if s.stopped {
		return nil
	}
	s.stopped = true
	if s.cancel != nil {
		s.cancel()
	}
	return nil
}

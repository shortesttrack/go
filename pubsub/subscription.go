package pubsub

import (
	"context"
	"cloud.google.com/go/pubsub"
	"st-go/errors"
	"sync"
)

type SubscriptionOpts struct {
	ProjectID string
	Name      string
}

type subscription struct {
	*pubsub.Client
	*pubsub.Subscription

	name string

	ctx    context.Context

	subscriptionContext context.Context
	subscriptionCancel func()
	err error
	stopMu sync.Mutex
	stopped bool
}

func NewSubscription(ctx context.Context, options SubscriptionOpts) (Subscription, error) {
	client, err := pubsub.NewClient(ctx, options.ProjectID)
	if err != nil {
		return nil, err
	}
	s := client.Subscription(options.Name)
	exists, err := s.Exists(ctx)
	if err != nil {
		client.Close()
		return nil, errors.New("invalid subscription: %s" + options.Name)
	}
	if !exists {
		return nil, errors.New("subscription " + options.Name + "does not exist")
	}
	return &subscription{
		Client: client,
		Subscription: s,
		name: options.Name,
		ctx: ctx,
		}, nil
}

func (s *subscription) Start() <-chan Message {
	output := make(chan Message)
	s.subscriptionContext, s.subscriptionCancel = context.WithCancel(s.ctx)
	go func(s *subscription, output chan Message) {
		defer close(output)
		err := s.Subscription.Receive(s.subscriptionContext, func(ctx context.Context, msg *pubsub.Message) {
			output <- &message{msg}
		})
		s.err = err
		s.Stop()
	}(s, output)
	return output
}

func (s *subscription) Err() error {
	return s.err
}

func (s *subscription) Stop() {
	s.stopMu.Lock()
	defer s.stopMu.Unlock()
	if s.stopped {
		return
	}
	s.stopped = true
	if s.subscriptionCancel != nil {
		s.subscriptionCancel()
	}
}

func (s *subscription) Context() context.Context {
	if s.subscriptionContext == nil {
		return s.ctx
	}
	return s.subscriptionContext
}
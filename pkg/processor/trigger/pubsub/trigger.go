/*
Copyright 2017 The Nuclio Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pubsub

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/nuclio/nuclio/pkg/common"
	"github.com/nuclio/nuclio/pkg/errors"
	"github.com/nuclio/nuclio/pkg/functionconfig"
	"github.com/nuclio/nuclio/pkg/processor/trigger"
	"github.com/nuclio/nuclio/pkg/processor/worker"

	pubsubClient "cloud.google.com/go/pubsub"
	"github.com/nuclio/logger"
	"github.com/rs/xid"
)

type pubsub struct {
	trigger.AbstractTrigger
	configuration      *Configuration
	stop               chan bool
	subscriptions      []*pubsubClient.Subscription
	client 				*pubsubClient.Client
}

func newTrigger(parentLogger logger.Logger,
	workerAllocator worker.Allocator,
	configuration *Configuration) (trigger.Trigger, error) {

	newTrigger := &pubsub{
		AbstractTrigger: trigger.AbstractTrigger{
			ID:              configuration.ID,
			Logger:          parentLogger.GetChild(configuration.ID),
			WorkerAllocator: workerAllocator,
			Class:           "async",
			Kind:            "pubsub",
		},
		configuration: configuration,
		stop:          make(chan bool),
	}

	return newTrigger, nil
}

func (p *pubsub) Start(checkpoint functionconfig.Checkpoint) error {
	var err error

	// TODO: find a better way to do this
	serviceAccountFilePath := "/tmp/service-account.json"
	ioutil.WriteFile(serviceAccountFilePath, []byte(p.configuration.Credentials.Contents), 0600)
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", serviceAccountFilePath)

	p.Logger.InfoWith("Starting",
		"subscriptions", p.configuration.Subscriptions,
		"projectID", p.configuration.ProjectID,
	)

	// pubsub client consumes namespace/project string to be created
	p.client, err = pubsubClient.NewClient(context.TODO(), p.configuration.ProjectID)
	if err != nil {
		return errors.Wrapf(err, "Can't connect to pubsub project")
	}

	for _, subscription := range p.configuration.Subscriptions {
		subscription := subscription

		go func() {
			err := p.receiveFromSubscription(&subscription)
			if err != nil {
				p.Logger.WarnWith("Failed to create subscription",
					"err", errors.GetErrorStackString(err, 10),
					"subscription", subscription)
			}
		}()
	}


	return nil
}

func (p *pubsub) Stop(force bool) (functionconfig.Checkpoint, error) {
	return nil, nil
}

func (p *pubsub) GetConfig() map[string]interface{} {
	return common.StructureToMap(p.configuration)
}

func (p *pubsub) receiveFromSubscription(subscriptionConfig *Subscription) error {
	ctx := context.TODO()

	// get subscription name
	subscriptionID := p.getSubscriptionID(subscriptionConfig)

	// get ack timeout
	ackDeadline, err := p.getAckDeadline(subscriptionConfig)
	if err != nil {
		return errors.Wrap(err, "Failed to parse ack deadline")
	}

	// try to create a subscription
	subscription, err := p.client.CreateSubscription(ctx, subscriptionID, pubsubClient.SubscriptionConfig{
		Topic: p.client.Topic(subscriptionConfig.Topic),
		AckDeadline: ackDeadline,
	})

	if err != nil {

		if !subscriptionConfig.Shared {
			return errors.Wrap(err, "Failed to create subscirption")
		}

		// try to use a subscription
		subscription = p.client.Subscription(subscriptionID)
	}

	// https://godoc.org/cloud.google.com/go/pubsub#ReceiveSettings
	subscription.ReceiveSettings.NumGoroutines = subscriptionConfig.MaxNumWorkers

	// create a channel of events
	eventsChan := make(chan *Event, subscriptionConfig.MaxNumWorkers)

	for eventIdx := 0; eventIdx < subscriptionConfig.MaxNumWorkers; eventIdx++ {
		eventsChan <- &Event{
			topic: subscriptionConfig.Topic,
		}
	}

	// listen to subscribed topic messages
	err = subscription.Receive(ctx, func(ctx context.Context, message *pubsubClient.Message) {

		// get an event
		event := <- eventsChan

		// set the message
		event.message = message

		// process the event, don't really do anything with response
		_, submitError, processError := p.AllocateWorkerAndSubmitEvent(event, p.Logger, 10 * time.Second)
		if submitError != nil {
			p.Logger.ErrorWith("Can't submit event", "error", submitError)

			message.Nack() // necessary to call on fail
			eventsChan <- event
			return
		}
		if processError != nil {
			p.Logger.ErrorWith("Can't process event", "error", processError)

			message.Nack()
			eventsChan <- event
			return
		}

		message.Ack()

		// return event to pool
		eventsChan <- event
	})

	if err != context.Canceled {
		return errors.Wrapf(err, "Message receiver cancelled")
	}

	return nil
}

func (p *pubsub) getSubscriptionID(subscriptionConfig *Subscription) string {
	subscriptionID := "nuclio-sub-" + subscriptionConfig.Topic

	// if multiple replicas share this subscription it must be named the same
	if subscriptionConfig.Shared {
		return subscriptionID
	}

	// if it's not shared, we must add a unique variable
	return subscriptionID + "-" + xid.New().String()
}

func (p *pubsub) getAckDeadline(subscriptionConfig *Subscription) (time.Duration, error) {
	var ackDeadlineString string

	if subscriptionConfig.AckDeadline != "" {
		ackDeadlineString = subscriptionConfig.AckDeadline
	} else if p.configuration.AckDeadline != "" {
		ackDeadlineString = p.configuration.AckDeadline
	} else {
		ackDeadlineString = "10s"
	}

	return time.ParseDuration(ackDeadlineString)
}

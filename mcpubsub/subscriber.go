package main

import (
	"bufio"
	"errors"
	"fmt"
)

type subscriberApi interface {
	Subscribe(topic string, callback Func) error
	Unsubscribe(topic string) error
}

type subscriber struct {
	id          int
	rw          *bufio.ReadWriter
	callBackMap map[string]Func
}

var (
	subscriberId                             = 1
	subscriptionMap map[string][]*subscriber = make(map[string][]*subscriber)
)

func newSubscriber(rw *bufio.ReadWriter) subscriberApi {
	subscriberObj := &subscriber{
		subscriberId,
		rw,
		make(map[string]Func),
	}
	subscriberId++
	return subscriberObj
}

func (subscriberObj *subscriber) handleMessage(message serverMessage) {
	callback := subscriberObj.callBackMap[message.Topic]
	callback(message.Message)
}

func (subscriberObj *subscriber) Subscribe(topic string, callback Func) error {
	subscriberObj.callBackMap[topic] = callback

	_, ok := subscriptionMap[topic]
	if !ok {
		subscriptionMap[topic] = make([]*subscriber, 0)
	}

	if hasElement(subscriptionMap[topic], *subscriberObj) {
		return errors.New("Already subscriber to " + topic)
	}

	subscriptionMap[topic] = append(subscriptionMap[topic], subscriberObj)
	fmt.Printf("Subscriber %d subscibed to %s \n", subscriberObj.id, topic)

	return nil
}

func (subscriberObj *subscriber) Unsubscribe(topic string) error {

	if !hasElement(subscriptionMap[topic], *subscriberObj) {
		return errors.New("No subscription exists.")
	}

	deleteStream(subscriptionMap[topic], subscriberObj)

	// Delete callback function to the topic
	delete(subscriberObj.callBackMap, topic)
	return nil
}

func deleteStream(subscribers []*subscriber, subscriber *subscriber) {
	for i := range subscribers {
		if subscribers[i].id == subscriber.id {
			subscribers = append(subscribers[:i], subscribers[i+1:]...)
		}
	}
}

func hasElement(subscriptions []*subscriber, subscriberObj subscriber) bool {
	for i := range subscriptions {
		if subscriptions[i].id == subscriberObj.id {
			return true
		}
	}
	return false
}

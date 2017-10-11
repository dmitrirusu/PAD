package main

import (
	"bufio"
)

type publisherApi interface {
	Publish(topic string, message string) error
}

type publisher struct {
	id int
	rw *bufio.ReadWriter
}

var publisherId = 1

func newPublisher(rw *bufio.ReadWriter) publisherApi {
	publisherObj := publisher{publisherId, rw, }
	publisherId++

	sendMessage(publisherObj.rw, serverMessage{
		Type: addPublisher},
	)
	return publisherObj
}

func (publisherObj publisher) Publish(topic string, message string) error {
	err := sendMessage(publisherObj.rw, serverMessage{
		Id:      publisherObj.id,
		Type:    published,
		Topic:   topic,
		Message: message,
	})
	return err
}

package mcpubsub

import (
	"bufio"
	"net"
)

type publisherApi interface {
	Publish(topic string, message string) error
}

type publisher struct {
	id int
	rw *bufio.ReadWriter
}

func newRW() (*bufio.ReadWriter, error) {

	conn, err := net.Dial("tcp", "localhost"+port)
	if err != nil {
		return nil, err
	}
	rw := bufio.NewReadWriter(bufio.NewReader(conn),
		bufio.NewWriter(conn))
	return rw, nil
}

func NewPub() publisherApi {
	rw, _ := newRW()
	publisherObj := publisher{
		clientId,
		rw,
	}

	clientId++
	return publisherObj
}

func (publisherObj publisher) Publish(topic string, message string) error {
	err := sendMessage(publisherObj.rw, serverMessage{
		Id:      publisherObj.id,
		Type:    messagePublished,
		Topic:   topic,
		Message: message,
	})
	return err
}

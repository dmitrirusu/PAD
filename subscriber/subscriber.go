package main

import (
	"bufio"
	"net"
	"encoding/gob"
	"log"
	"fmt"
)

const (
	newSubscriber    = "New subscriber"
)

var port = ":12449"

type subscriberApi interface {
	Subscribe(topic string, callback Func) error
}

type Func func(string)

type subscriber struct {
	rw          *bufio.ReadWriter
	callBackMap map[string]Func
}

type serverMessage struct {
	Id      int
	Type    string
	Topic   string
	Message string
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

func (sub *subscriber) Subscribe(topic string, callback Func) error {
	sub.callBackMap[topic] = callback

	message := serverMessage{
		Type:  newSubscriber,
		Topic: topic,
	}

	enc := gob.NewEncoder(sub.rw)
	err := enc.Encode(&message)
	if err != nil {
		log.Println("Error while sending message.", message, err)
		return err
	}
	sub.rw.Flush()

	for {
		msg, err := receiveMessage(sub.rw)
		if err != nil {
			fmt.Println("Could not subscribe ")
			break
		}
		sub.callBackMap[topic](msg.Message)
	}

	return nil
}

func receiveMessage(rw *bufio.ReadWriter) (message serverMessage, err error) {
	dec := gob.NewDecoder(rw)
	err = dec.Decode(&message)
	if err != nil {
		log.Println("Error decoding data", err)
		return message, err
	}
	rw.Flush()
	return message, nil
}

func NewSub() subscriberApi {
	rw, _ := newRW()
	subscriberObj := &subscriber{
		rw,
		make(map[string]Func),
	}

	return subscriberObj
}

func (sub *subscriber) hasCallback(message serverMessage) bool {
	_, ok := sub.callBackMap[message.Topic]
	return ok
}

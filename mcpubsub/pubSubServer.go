package main

import (
	"bufio"
	"errors"
	"log"
	"net"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

var (
	PORT                        = ":14610"
	streams []*bufio.ReadWriter = make([]*bufio.ReadWriter, 0)
)

type pubSubServerApi interface {
	newRW() (*bufio.ReadWriter, error)
	start() error
}

type pubSubServer struct {
}

type serverMessage struct {
	Id      int
	Type    string
	Topic   string
	Message string
}

func (message *serverMessage) toJson() []byte {
	marshal, _ := json.Marshal(message)
	return marshal
}

func sendMessage(rw *bufio.ReadWriter, message serverMessage) error {
	rw.Write(message.toJson())
	rw.Flush()
	fmt.Println("sent message " + message.Message)
	return nil
}

func receiveMessage(rw *bufio.ReadWriter) (message serverMessage, err error) {
	fmt.Println("Receive")
	var msg serverMessage
	all, _ := ioutil.ReadAll(rw)
	json.Unmarshal(all, &msg)

	fmt.Println("Received message " + msg.Message)
	return msg, nil
}

func (server pubSubServer) newRW() (*bufio.ReadWriter, error) {
	conn, err := net.Dial("tcp", "localhost"+PORT)
	if err != nil {
		return nil, err
	}
	log.Println("New client added.")
	rw := bufio.NewReadWriter(bufio.NewReader(conn),
		bufio.NewWriter(conn))
	return rw, nil
}

func (server pubSubServer) start() error {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		return errors.New("Not able to accept connections.")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}
		log.Println("Connection Accepted.")

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	for {
		message, err := receiveMessage(rw)
		if err != nil {
			log.Println("error in handleConnection", err)
			return
		}
		handleRequestMessage(message, rw)
	}
}

func handleRequestMessage(message serverMessage, rw *bufio.ReadWriter) {
	var err error = nil
	switch message.Type {
	case addClientMessage:
		streams = append(streams, rw)
	case addPublisher:
		message.Id = publisherId
		publisherId++
		err = sendMessage(rw, message)
	case published:
		for idx, stream := range streams {
			log.Println(idx, stream)

			fmt.Println("Publishing msg")
			err = sendMessage(stream, message)
			if err != nil {
				log.Println("Bad stream can't write.")
			}
		}
	case subscribed:
		return
	case unsubsrcibed:
		return
	default:
		log.Println("Unrecognised message.")
	}
}

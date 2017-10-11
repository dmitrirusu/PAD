package mcpubsub

import (
	"bufio"
	"encoding/gob"
	"errors"
	"log"
	"net"
	"fmt"
	"io/ioutil"
	"time"
	"encoding/json"
)

type broker interface {
	newRW() (*bufio.ReadWriter, error)
	Start() error
}

type serverMessage struct {
	Id      int
	Type    string
	Topic   string
	Message string
}

func sendMessage(rw *bufio.ReadWriter, message serverMessage) error {
	enc := gob.NewEncoder(rw)
	err := enc.Encode(&message)
	if err != nil {
		log.Println("Error while sending message.", message, err)
		return err
	}
	rw.Flush()
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

func Start() error {
	listen, err := net.Listen("tcp", "localhost"+port)
	if err != nil {
		return errors.New("Not able to accept connections.")
	}

	restoreData()
	go broadCaster()
	go startBackupTimer()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}
		log.Println("Connection Accepted.")

		go handleConnection(conn)
	}
}

func broadCaster() {
	for {
		for _, queue := range pipeline {
			message := queue.Pop()
			broadCastMessage(message)
		}
	}
}

func broadCastMessage(message serverMessage) (err error) {

	if message.Id == 0 {
		return errors.New("Invalid Message, Id is null.")
	}

	streams, ok := subscriptionMap[message.Topic]
	log.Println(message, streams, ok)
	if !ok {
		return nil
	}

	for _, subscriber := range subscriptionMap[message.Topic] {
		sendMessage(subscriber, message)
	}
	return nil
}

func startBackupTimer() {
	fmt.Println("Start backup timer")

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				backup()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func backup() {
	bytes, _ := json.Marshal(pipeline)
	fmt.Println("Backup ", string(bytes))
	ioutil.WriteFile(filePath, bytes, 0644)
}

func restoreData() {
	time.Sleep(time.Second * 2)

	var restore map[string]*Queue

	file, _ := ioutil.ReadFile(filePath)
	json.Unmarshal(file, &restore)
	for topic, queue := range restore {
		fmt.Println("Restoring topic" + topic)
		_, ok := pipeline[topic]
		if !ok {
			pipeline[topic] = NewQueue()
		}
		for _, msg := range queue.Messages {
			fmt.Println("Restoring message " + msg.Message)
			pipeline[topic].Push(msg)
		}
	}
	ioutil.WriteFile(filePath, []byte(""), 0644)
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
		handleMessage(message, rw)
	}
}

func handleMessage(message serverMessage, rw *bufio.ReadWriter) {
	fmt.Println("Received message type " + message.Type)
	switch message.Type {
	case newPublisher:
		fmt.Printf("Publisher connected")
	case messagePublished:
		_, ok := pipeline[message.Topic]
		if !ok {
			pipeline[message.Topic] = NewQueue()
		}
		pipeline[message.Topic].Push(message)
	case newSubscriber:
		_, ok := subscriptionMap[message.Topic]
		if !ok {
			sendMessage(rw, serverMessage{
				Type:    subError,
				Message: "Could not subscribe",
			})

		} else {
			subscriptionMap[message.Topic] = append(subscriptionMap[message.Topic], rw)
			fmt.Println("Subsacribed to " + message.Topic)
		}

	default:
		sendMessage(rw, serverMessage{
			Type:    unknownTypeError,
			Message: unknownTypeError,
		})
		log.Println("Unrecognised message.")
	}
}

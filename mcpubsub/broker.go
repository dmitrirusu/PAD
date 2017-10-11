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

	go broadCaster()
	go startBackupTimer()
	go restoreData()

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
	ioutil.WriteFile(filePath, bytes, 0644)
}

func restoreData() {
	time.Sleep(time.Second * 2)

	var messages []serverMessage

	file, _ := ioutil.ReadFile(filePath)
	json.Unmarshal(file, &messages)
	for _, msg := range messages {
		fmt.Println("Restoring " + msg.Message)
		_, ok := pipeline[msg.Topic]
		if !ok {
			pipeline[msg.Topic] = NewQueue()
		}
		pipeline[msg.Topic].Push(msg)
	}

	backup()
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
	fmt.Println("Received message type" + message.Type)
	switch message.Type {
	case newPublisher:
		fmt.Printf("Publisher connected")
	case messagePublished:
		_, ok := pipeline[message.Topic]
		if !ok {
			pipeline[message.Topic] = NewQueue()
		}
		pipeline[message.Topic].Push(message)
		subscriptionMap[message.Topic] = make([]*bufio.ReadWriter, 0)
	case newSubscriber:
		_, ok := subscriptionMap[message.Topic]
		if !ok {
			rw.Write([]byte("Unrecognized topic"))
		} else {
			subscriptionMap[message.Topic] = append(subscriptionMap[message.Topic], rw)
			fmt.Println("Subsacribed to " + message.Topic)
		}
	default:
		rw.Write([]byte("Unrecognised message"))
	}
}

package main

import (
	"bufio"
	"errors"
	"log"
	"time"
	"fmt"
	"encoding/gob"
	"os"
)

type Func func(string)

type pubSubApi interface {
	NewPublisher() publisherApi
	NewSubscriber() subscriberApi
}

type pubSubFactory struct {
	server pubSubServerApi
	rw     *bufio.ReadWriter
}

func PubSubApi() (pubSubApi, error) {
	initChannelMap()
	//go restoreData()
	//startBackuptTimer()

	server := pubSubServer{}
	rw, err := server.newRW()
	if err != nil {
		log.Println("Error in starting service", err)
		return nil, err
	}

	go listener(rw)
	go broadCaster()

	sendMessage(rw, serverMessage{
		Type: addClientMessage,
	})

	pubSubObj := pubSubFactory{
		server: server,
		rw:     rw,
	}

	return pubSubObj, nil
}

func listener(rw *bufio.ReadWriter) {
	for {
		fmt.Println("Pushed to chnnel")

		message, err := receiveMessage(rw)
		if err != nil {
			log.Println("Not able to receive messages.", err)
			return
		}

		channelMap[message.Type].Push(message)
	}
}

func startBackuptTimer() {
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
	file, _ := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	enc := gob.NewEncoder(file)
	enc.Encode(&channelMap[published].messages)
}

func restoreData() {
	//time.Sleep(time.Second * 2)
	//var messages []serverMessage
	//
	//file, _ := os.OpenFile(filePath, os.O_RDONLY, 0644)
	//dec := gob.NewDecoder(file)
	//dec.Decode(&messages)
	//
	//for _, msg := range messages {
	//	fmt.Println("Restored message : " + msg.Message)
	//	channelMap[msg.Type].Push(msg)
	//}
}

func broadCaster() {
	for {
		fmt.Println("Broadcaster " + string(channelMap[published].ToJson()))
		message := channelMap[published].Pop()
		broadCastMessage(message)
	}
}

func broadCastMessage(message serverMessage) (err error) {

	log.Println("Broadcasting message", message)

	if message.Id == 0 {
		return errors.New("Invalid Message, Id is null.")
	}

	streams, ok := subscriptionMap[message.Topic]
	log.Println(message, streams, ok)
	if !ok {
		return nil
	}

	for _, subscriber := range subscriptionMap[message.Topic] {
		fmt.Println(subscriber.id)
		subscriber.handleMessage(message)
	}
	return nil
}

func PubSubApiServerStart() (pubSubApi, error) {
	pubSubObj := pubSubFactory{
		server: pubSubServer{},
	}
	err := pubSubObj.server.start()
	if err != nil {
		return nil, err
	}
	return pubSubObj, nil
}

func (pubSubObj pubSubFactory) NewPublisher() publisherApi {
	return newPublisher(pubSubObj.rw)
}

func (pubSubObj pubSubFactory) NewSubscriber() subscriberApi {
	return newSubscriber(pubSubObj.rw)
}

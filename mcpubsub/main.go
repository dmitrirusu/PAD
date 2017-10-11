package main

import (
	"fmt"
	"log"
	"time"
)

func main() {

	fmt.Println("Starting...")
	go PubSubApiServerStart()
	fmt.Println("Server started")

	pubSubApi, err := PubSubApi()
	if err != nil {
		log.Println("Error in listening.")
		return
	}
	publisher := pubSubApi.NewPublisher()
	subscriber := pubSubApi.NewSubscriber()

	subscriber.Subscribe("got", func(message string) {
		fmt.Println("Game of thrones: ", message)
	})
	publisher.Publish("got", "All hail Lord Baelish")
	publisher.Publish("got", "test")

	subscriber.Subscribe("rm", func(message string) {
		fmt.Println("Rick and Morty", message)
	})

	time.Sleep(time.Millisecond)
	subscriber.Unsubscribe("got")
	publisher.Publish("rm", "Dabba wabba doo")
	publisher.Publish("got", "All hail Lord Baelish")

	time.Sleep(time.Second * 50)

}

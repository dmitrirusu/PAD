package mcpubsub

import (
	"fmt"
	"time"
)

func main() {

	fmt.Println("Starting...")
	go Start()
	fmt.Println("Server started")

	time.Sleep(time.Second * 5)

	//publisher := NewPub()
	////subscriber := NewSub()
	////
	////subscriber.Subscribe("got", func(message string) {
	////	fmt.Println("Game of thrones: ", message)
	////})
	//publisher.Publish("got", "All hail Lord Baelish")
	//publisher.Publish("got", "test")
	//
	////subscriber.Subscribe("rm", func(message string) {
	////	fmt.Println("Rick and Morty", message)
	////})
	//
	//time.Sleep(time.Millisecond)
	////subscriber.UnSubscribe("got")
	////publisher.Publish("rm", "Dabba wabba doo")
	//publisher.Publish("got", "All hail Lord Baelish")

}

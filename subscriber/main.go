package main

import (
	"bufio"
	"os"
	"fmt"
	"time"
)

func main() {

	sub := NewSub()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Topic: ")

	topic, _ := reader.ReadString('\n')

	go sub.Subscribe(topic, func(message string) {
		fmt.Println("Received message: ", message)
	})
	time.Sleep(time.Second * 10000000)
}

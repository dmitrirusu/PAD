package main

import (
	"github.com/mc-lovin_work/mcpubsub"
	"bufio"
	"os"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"time"
)

func main() {

	pub := mcpubsub.NewPub()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Topic: ")
	topic, _ := reader.ReadString('\n')

	for {


		time.Sleep(time.Second * 1)
		u, _ := uuid.NewV4()

		pub.Publish(topic, u.String())
	}
}

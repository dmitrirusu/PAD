package main

import ("bufio"
	"os"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"time"
	"github.com/git_pad/mcpubsub"
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

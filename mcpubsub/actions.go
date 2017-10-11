package main

const (
	addClientMessage = "CLIENT ADDED"
	addPublisher     = "PUBLISHER ADDED"
	published        = "PUBLISHER PUBLISHED"
	subscribed       = "SUBSCRIBER SUBSCRIBED"
	unsubsrcibed     = "SUBSCRIBER UNSUBSCRIBED"
	filePath         = "/home/sylar/backup"
)

var (
	channelMap map[string]*Queue = make(map[string]*Queue)
)

func initChannelMap() {
	channelMap[addPublisher] = NewQueue()
	channelMap[published] = NewQueue()
	channelMap[subscribed] = NewQueue()
	channelMap[unsubsrcibed] = NewQueue()
}

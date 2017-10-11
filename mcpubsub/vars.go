package mcpubsub

import "bufio"

const (
	newPublisher     = "New publisher"
	messagePublished = "Message published"
	newSubscriber    = "New subscriber"
	filePath         = "/home/sylar/backup.json"
)

var (
	subscriptionMap map[string][]*bufio.ReadWriter = make(map[string][]*bufio.ReadWriter)
	pipeline        map[string]*Queue              = make(map[string]*Queue)
	port                                           = ":12449"
	clientId                                       = 1
)

package mcpubsub

import "bufio"

const (
	connect          = "Client connected"
	newPublisher     = "New publisher"
	subError         = "Error"
	unknownTypeError = "Unknown message type"
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

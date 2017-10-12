package mcpubsub

import (
	"sync"
)

type Queue struct {
	Messages []serverMessage `json:"messages"`
	mutex    *sync.Cond
}

func (q *Queue) Push(data serverMessage) {
	q.mutex.L.Lock()
	defer q.mutex.L.Unlock()
	q.Messages = append(q.Messages, data)
}

func NewQueue() *Queue {
	return &Queue{
		Messages: make([]serverMessage, 0),
		mutex:    sync.NewCond(new(sync.Mutex)),
	}
}

func (q *Queue) Pop() serverMessage {
	q.mutex.L.Lock()
	defer q.mutex.L.Unlock()

	if (len(q.Messages) > 0) {
		msg := q.Messages[len(q.Messages)-1]
		q.Messages = q.Messages[:len(q.Messages)-1]
		return msg
	} else {
		return serverMessage{Type:"nil"}
	}
}

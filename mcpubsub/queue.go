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
	q.mutex.Signal()
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

	for len(q.Messages) == 0 {
		q.mutex.Wait()
	}

	msg := q.Messages[len(q.Messages)-1]
	q.Messages = q.Messages[:len(q.Messages)-1]
	return msg
}

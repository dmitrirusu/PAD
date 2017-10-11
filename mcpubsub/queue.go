package mcpubsub

import (
	"sync"
)

type Queue struct {
	messages []serverMessage
	mutex    *sync.Cond
}

func (q *Queue) Push(data serverMessage) {
	q.mutex.L.Lock()
	defer q.mutex.L.Unlock()
	q.messages = append(q.messages, data)
	q.mutex.Signal()
}

func NewQueue() *Queue {
	return &Queue{
		messages: make([]serverMessage, 0),
		mutex:    sync.NewCond(new(sync.Mutex)),
	}
}

func (q *Queue) Pop() serverMessage {
	q.mutex.L.Lock()
	defer q.mutex.L.Unlock()

	for len(q.messages) == 0 {
		q.mutex.Wait()
	}

	msg := q.messages[len(q.messages)-1]
	q.messages = q.messages[:len(q.messages)-1]
	return msg
}

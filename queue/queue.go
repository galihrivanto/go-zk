package queue

import (
	"errors"
	"sync"
)

// define common queue error
var (
	ErrQueueEmpty    = errors.New("queue is empty")
	ErrInternalError = errors.New("internal queue error")
)

// Default queue
var Default = NewQueue()

// Queuer represent simple FIFO queue
// queue will divided by queue name
// eg:
// Queue 1: [item1, item2]
// Queue 2: [item3, item4]
type Queuer interface {

	// Push add new item in queue
	Push(queueName string, item string) error

	// Pop get oldest item in queue
	// if no item found then return nil + ErrQueueEmpty
	Pop(queueName string) (string, error)

	// Len request queue length
	Len(queueName string) int
}

// InMemoryQueue implement Queuer using memory
// to store queue item
type InMemoryQueue struct {
	// queue map
	queues sync.Map
}

func (q *InMemoryQueue) getQueue(queueName string) *internalQueue {
	// retrive queue by name
	// if not exist then create a new one
	v, _ := q.queues.LoadOrStore(queueName, &internalQueue{})

	if queue, ok := v.(*internalQueue); ok {
		return queue
	}

	return nil
}

// Push implements Queuer.Push
func (q *InMemoryQueue) Push(queueName string, item string) error {
	queue := q.getQueue(queueName)
	if queue == nil {
		return ErrInternalError
	}

	// push item to internal queue
	return queue.push(item)
}

// Pop implements Queuer.Pop
func (q *InMemoryQueue) Pop(queueName string) (string, error) {
	queue := q.getQueue(queueName)
	if queue == nil {
		return "", ErrInternalError
	}

	// pop item to internal queue
	return queue.pop()
}

// Len implements Queuer.Len
func (q *InMemoryQueue) Len(queueName string) int {
	queue := q.getQueue(queueName)
	if queue == nil {
		return 0
	}

	return queue.len()
}

// NewQueue create in memory queue
func NewQueue() Queuer {
	return &InMemoryQueue{}
}

// internalQueue is helper around FIFO queue
type internalQueue struct {
	sync.Mutex

	// first item is oldest
	// last item is newest
	items []string
}

func (q *internalQueue) push(item string) error {
	q.Lock()
	defer q.Unlock()

	if q.items == nil {
		q.items = make([]string, 0)
	}

	q.items = append(q.items, item)

	return nil
}

func (q *internalQueue) pop() (string, error) {
	q.Lock()
	defer q.Unlock()

	if q.items == nil {
		return "", ErrQueueEmpty
	}

	// pop items
	item := q.items[0]
	q.items = q.items[1:]

	return item, nil
}

func (q *internalQueue) len() int {
	q.Lock()
	defer q.Unlock()

	return len(q.items)
}

// Push add new item in queue
func Push(queueName string, item string) error {
	return Default.Push(queueName, item)
}

// Pop get oldest item in queue
func Pop(queueName string) (string, error) {
	return Default.Pop(queueName)
}

// Len get queue length
func Len(queueName string) int {
	return Default.Len(queueName)
}

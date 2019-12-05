package queue

import "testing"

// In memory queue test
func TestQueue(t *testing.T) {
	queue := NewQueue()
	queueName := "TEST"

	// drain old queue'ed item first
	for queue.Len(queueName) > 0 {
		queue.Pop(queueName)
	}

	// put first item
	queue.Push(queueName, "ITEM1")

	if queue.Len(queueName) != 1 {
		t.Error("Expected length 1")
		t.FailNow()
	}

	// put second item
	queue.Push(queueName, "ITEM2")

	if queue.Len(queueName) != 2 {
		t.Error("Expected length 2")
		t.FailNow()
	}

	// put third item
	queue.Push(queueName, "ITEM3")

	if queue.Len(queueName) != 3 {
		t.Error("Expected length 3")
		t.FailNow()
	}

	// pop queue
	item, _ := queue.Pop(queueName)

	if item != "ITEM1" {
		t.Error("Expected item ITEM1 but returned", item)
	}
}

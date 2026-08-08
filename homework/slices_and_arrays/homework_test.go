package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type circularQueue struct {
	values   []int
	head     int
	tail     int
	count    int
	capacity int
}

func NewCircularQueue(size int) circularQueue {
	return circularQueue{
		values:   make([]int, size),
		head:     0,
		tail:     0,
		count:    0,
		capacity: size,
	}
}

func (q *circularQueue) Push(value int) bool {
	if q.Full() {
		return false
	}
	q.values[q.tail] = value
	q.tail = (q.tail + 1) % q.capacity
	q.count++
	return true
}

func (q *circularQueue) Pop() bool {
	if q.Empty() {
		return false
	}
	q.head = (q.head + 1) % q.capacity
	q.count--
	return true
}

func (q *circularQueue) Front() int {
	if q.Empty() {
		return -1
	}
	return q.values[q.head]
}

func (q *circularQueue) Back() int {
	if q.Empty() {
		return -1
	}
	// Так как tail указывает на следующий пустой слот,
	// последний элемент находится по индексу (tail - 1 + capacity) % capacity
	lastIndex := (q.tail - 1 + q.capacity) % q.capacity
	return q.values[lastIndex]
}

func (q *circularQueue) Empty() bool {
	return q.count == 0
}

func (q *circularQueue) Full() bool {
	return q.count == q.capacity
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue(queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}

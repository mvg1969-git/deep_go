package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Task struct {
	Identifier int
	Priority   int
}

type scheduler struct {
	heap []*Task
	// Хранит связь: ID задачи -> её текущий индекс в массиве heap
	index map[int]int
}

func NewScheduler() scheduler {
	return scheduler{
		heap:  make([]*Task, 0),
		index: make(map[int]int),
	}
}

func (s *scheduler) AddTask(Task *Task) {
	s.heap = append(s.heap, Task)
	taskIndex := len(s.heap) - 1
	s.index[Task.Identifier] = taskIndex
	s.up(taskIndex)
}

func (s *scheduler) ChangeTaskPriority(taskid int, newpriority int) {
	idx, exists := s.index[taskid]
	if !exists {
		return
	}

	oldPriority := s.heap[idx].Priority
	s.heap[idx].Priority = newpriority

	// Если приоритет вырос, двигаем вверх, иначе — вниз (Max-Heap)
	if newpriority > oldPriority {
		s.up(idx)
	} else if newpriority < oldPriority {
		s.down(idx)
	}
}

func (s *scheduler) GetTask() *Task {
	if len(s.heap) == 0 {
		return &Task{}
	}

	root := s.heap[0]
	lastIdx := len(s.heap) - 1

	// Переносим последний элемент на место корня
	s.swap(0, lastIdx)
	s.heap = s.heap[:lastIdx]
	delete(s.index, root.Identifier)

	if len(s.heap) > 0 {
		s.down(0)
	}

	return root
}

// Вспомогательные методы для работы с Max-Heap

func (s *scheduler) swap(i, j int) {
	s.heap[i], s.heap[j] = s.heap[j], s.heap[i]
	s.index[s.heap[i].Identifier] = i
	s.index[s.heap[j].Identifier] = j
}

func (s *scheduler) up(i int) {
	for {
		parent := (i - 1) / 2
		if i == 0 || s.heap[parent].Priority >= s.heap[i].Priority {
			break
		}
		s.swap(parent, i)
		i = parent
	}
}

func (s *scheduler) down(i int) {
	for {
		left := 2*i + 1
		if left >= len(s.heap) || left < 0 {
			break
		}

		// Ищем наибольшего потомка
		largest := left
		if right := left + 1; right < len(s.heap) && s.heap[right].Priority > s.heap[left].Priority {
			largest = right
		}

		if s.heap[i].Priority >= s.heap[largest].Priority {
			break
		}
		s.swap(i, largest)
		i = largest
	}
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(&task1)
	scheduler.AddTask(&task2)
	scheduler.AddTask(&task3)
	scheduler.AddTask(&task4)
	scheduler.AddTask(&task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, *task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, *task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	assert.Equal(t, task1, *task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, *task)
}

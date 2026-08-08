package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Trace(stacks [][]uintptr) []uintptr {
	var result []uintptr
	var visited = make(map[uintptr]bool)
	var onStack = make(map[uintptr]bool)

	// Предварительный этап: фиксируем все адреса, которые явно переданы на стеках
	for _, stack := range stacks {
		for _, ptr := range stack {
			if ptr != 0 {
				onStack[ptr] = true
			}
		}
	}

	// Основной этап: обход каждого стека по отдельности
	for _, stack := range stacks {
		var queue []uintptr

		// 1. Собираем все уникальные корни текущего стека
		for _, ptr := range stack {
			if ptr != 0 && !visited[ptr] {
				visited[ptr] = true
				result = append(result, ptr)
				queue = append(queue, ptr)
			}
		}

		// 2. Раскрываем их (BFS в рамках текущего стека)
		for len(queue) > 0 {
			curr := queue[0]
			queue = queue[1:]

			// Читаем значение из памяти
			nextPtr := *(*uintptr)(unsafe.Pointer(curr))

			// Добавляем объект, ТОЛЬКО если его нет на других стеках и мы его еще не видели
			if nextPtr != 0 && !visited[nextPtr] && !onStack[nextPtr] {
				visited[nextPtr] = true
				result = append(result, nextPtr)
				queue = append(queue, nextPtr)
			}
		}
	}

	return result
}

func TestTrace(t *testing.T) {
	var heapObjects = []int{
		0x00, 0x00, 0x00, 0x00, 0x00,
	}

	var heapPointer1 *int = &heapObjects[1]
	var heapPointer2 *int = &heapObjects[2]
	var heapPointer3 *int = nil
	var heapPointer4 **int = &heapPointer3

	var stacks = [][]uintptr{
		{
			uintptr(unsafe.Pointer(&heapPointer1)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[0])),
			0x00, 0x00, 0x00, 0x00,
		},
		{
			uintptr(unsafe.Pointer(&heapPointer2)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[1])),
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[2])),
			uintptr(unsafe.Pointer(&heapPointer4)), 0x00, 0x00, 0x00,
		},
		{
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[3])),
		},
	}

	pointers := Trace(stacks)
	expectedPointers := []uintptr{
		uintptr(unsafe.Pointer(&heapPointer1)),
		uintptr(unsafe.Pointer(&heapObjects[0])),
		uintptr(unsafe.Pointer(&heapPointer2)),
		uintptr(unsafe.Pointer(&heapObjects[1])),
		uintptr(unsafe.Pointer(&heapObjects[2])),
		uintptr(unsafe.Pointer(&heapPointer4)),
		uintptr(unsafe.Pointer(&heapPointer3)),
		uintptr(unsafe.Pointer(&heapObjects[3])),
	}

	assert.True(t, reflect.DeepEqual(expectedPointers, pointers))
}

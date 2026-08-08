package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Defragment(memory []byte, pointers []unsafe.Pointer) {
	if len(pointers) == 0 {
		return
	}

	// Шаг 1: Переносим данные в начало и обновляем указатели
	writeIdx := 0
	for i := 0; i < len(pointers); i++ {
		// Читаем значение по старому адресу
		val := *(*byte)(pointers[i])

		// Записываем в целевую ячейку памяти
		memory[writeIdx] = val

		// Обновляем указатель, чтобы он смотрел на новый адрес
		pointers[i] = unsafe.Pointer(&memory[writeIdx])

		writeIdx++
	}

	// Шаг 2: Зануляем оставшуюся память в конце
	for i := writeIdx; i < len(memory); i++ {
		memory[i] = 0x00
	}
}

func TestDefragmentation(t *testing.T) {
	var fragmentedMemory = []byte{
		0xFF, 0x00, 0x00, 0x00,
		0x00, 0xFF, 0x00, 0x00,
		0x00, 0x00, 0xFF, 0x00,
		0x00, 0x00, 0x00, 0xFF,
	}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[5]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[15]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[3]),
	}

	var defragmentedMemory = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}

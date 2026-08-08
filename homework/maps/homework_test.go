package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

// Узел бинарного дерева поиска
type node struct {
	key   int
	value int
	left  *node
	right *node
}

// Структура упорядоченного словаря
type orderedMap struct {
	root  *node
	count int
}

// Конструктор
func NewOrderedMap() orderedMap {
	return orderedMap{}
}

// Вставка пары ключ-значение
func (m *orderedMap) Insert(key, value int) {
	m.root = m.insertNode(m.root, key, value)
}

func (m *orderedMap) insertNode(current *node, key, value int) *node {
	if current == nil {
		m.count++
		return &node{key: key, value: value}
	}
	if key < current.key {
		current.left = m.insertNode(current.left, key, value)
	} else if key > current.key {
		current.right = m.insertNode(current.right, key, value)
	} else {
		current.value = value // Обновление значения, если ключ уже есть
	}
	return current
}

// Удаление ключа
func (m *orderedMap) Erase(key int) {
	var deleted bool
	m.root, deleted = m.eraseNode(m.root, key)
	if deleted {
		m.count--
	}
}

func (m *orderedMap) eraseNode(current *node, key int) (*node, bool) {
	if current == nil {
		return nil, false
	}

	var deleted bool
	if key < current.key {
		current.left, deleted = m.eraseNode(current.left, key)
		return current, deleted
	}
	if key > current.key {
		current.right, deleted = m.eraseNode(current.right, key)
		return current, deleted
	}

	// Узел найден
	deleted = true

	// Случай 1: Нет детей или только один ребенок
	if current.left == nil {
		return current.right, deleted
	}
	if current.right == nil {
		return current.left, deleted
	}

	// Случай 2: Два ребенка. Ищем минимальный элемент в правом поддереве
	minRight := current.right
	for minRight.left != nil {
		minRight = minRight.left
	}

	current.key = minRight.key
	current.value = minRight.value

	// Удаляем дубликат минимального узла из правого поддерева
	current.right, _ = m.eraseNode(current.right, minRight.key)

	return current, deleted
}

// Проверка наличия ключа
func (m *orderedMap) Contains(key int) bool {
	current := m.root
	for current != nil {
		if key == current.key {
			return true
		} else if key < current.key {
			current = current.left
		} else {
			current = current.right
		}
	}
	return false
}

// Получение количества элементов
func (m *orderedMap) Size() int {
	return m.count
}

// Симметричный обход дерева (In-order traversal) для сохранения порядка
func (m *orderedMap) ForEach(action func(int, int)) {
	m.inOrder(m.root, action)
}

func (m *orderedMap) inOrder(current *node, action func(int, int)) {
	if current == nil {
		return
	}
	m.inOrder(current.left, action)
	action(current.key, current.value)
	m.inOrder(current.right, action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}

package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	notemptystruct bool
}

type MessageService struct {
	// not need to implement
	notemptystruct bool
}

// container хранит фабричные функции для создания объектов.
type container struct {
	factories map[string]func() interface{}
}

// newcontainer инициализирует контейнер и его внутреннюю карту.
func NewContainer() *container {
	return &container{
		factories: make(map[string]func() interface{}),
	}
}

// registertype сохраняет конструктор по заданному имени.
func (c *container) RegisterType(name string, constructor interface{}) {
	// Приводим интерфейс к конкретному типу функции, ожидаемому в тесте.
	if factory, ok := constructor.(func() interface{}); ok {
		c.factories[name] = factory
	}
}

// resolve вызывает конструктор и возвращает новый экземпляр или ошибку.
func (c *container) Resolve(name string) (interface{}, error) {
	factory, exists := c.factories[name]
	if !exists {
		return nil, fmt.Errorf("service %s not found", name)
	}

	// Вызов функции гарантирует создание нового объекта при каждом запросе
	return factory(), nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}

package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MultiError хранит список всех накопленных ошибок
type MultiError struct {
	Errors []error
}

// Error форматирует вывод в строгом соответствии с ожидаемым тестом сообщением
func (e *MultiError) Error() string {
	if e == nil || len(e.Errors) == 0 {
		return ""
	}

	var sb strings.Builder
	// Форматируем заголовок с количеством ошибок
	sb.WriteString(fmt.Sprintf("%d errors occured:\n", len(e.Errors)))

	// Собираем ошибки через знак табуляции и звездочку
	for _, err := range e.Errors {
		sb.WriteString(fmt.Sprintf("\t* %s", err.Error()))
	}

	// Добавляем финальный перенос строки, как в тесте
	sb.WriteString("\n")
	return sb.String()
}

// Append объединяет существующие ошибки в один MultiError
func Append(err error, errs ...error) *MultiError {
	result := &MultiError{}

	addErr := func(err error) {
		if me, ok := err.(*MultiError); ok {
			if me != nil {
				result.Errors = append(result.Errors, me.Errors...)
			}
		} else {
			result.Errors = append(result.Errors, err)
		}
	}

	// Если базовая ошибка уже является MultiError, копируем ее содержимое
	if err != nil {
		addErr(err)
	}

	// Добавляем все остальные переданные ошибки
	for _, e := range errs {
		addErr(e)
	}

	// Если ошибок не набралось, возвращаем nil
	if len(result.Errors) == 0 {
		return nil
	}

	return result
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}

package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Serialize(p Person) string {
	val := reflect.ValueOf(p)
	typ := reflect.TypeOf(p)
	var lines []string

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)

		// Получаем значение тега properties
		tag := structField.Tag.Get("properties")
		if tag == "" {
			continue
		}

		// Разделяем имя ключа и параметры (например, omitempty)
		parts := strings.Split(tag, ",")
		key := parts[0]
		omitempty := len(parts) > 1 && parts[1] == "omitempty"

		// Проверяем, пустое ли поле, для omitempty
		if omitempty && field.IsZero() {
			continue
		}

		// Форматируем базовые типы в строку
		var strVal string
		switch field.Kind() {
		case reflect.String:
			strVal = field.String()
		case reflect.Int:
			strVal = fmt.Sprintf("%d", field.Int())
		case reflect.Bool:
			strVal = fmt.Sprintf("%t", field.Bool())
		default:
			strVal = fmt.Sprintf("%v", field.Interface())
		}

		lines = append(lines, fmt.Sprintf("%s=%s", key, strVal))
	}

	return strings.Join(lines, "\n")
}

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}

package main

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	b := strings.Builder{}

	p := reflect.ValueOf(person)
	tpe := reflect.TypeOf(person)

	numField := tpe.NumField()

	for i := range numField {
		tags := strings.Split(tpe.Field(i).Tag.Get("properties"), ",")

		if len(tags) == 0 {
			continue
		}

		name := tags[0]
		value := p.Field(i)

		omitempty := len(tags) > 1 && slices.Contains(tags[1:], "omitempty")

		if value.IsZero() && omitempty {
			continue
		}

		_, _ = b.WriteString(fmt.Sprintf("%s=%v", name, value))

		if i != numField-1 {
			_, _ = b.WriteRune('\n')
		}
	}

	return b.String()
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

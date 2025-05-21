package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errs []error
}

func (e *MultiError) Error() string {
	if len(e.errs) == 0 {
		return ""
	}

	b := strings.Builder{}

	_, _ = b.WriteString(fmt.Sprintf("%d errors occured:\n", len(e.errs)))

	for _, err := range e.errs {
		_, _ = b.WriteString(fmt.Sprintf("\t* %s", err.Error()))
	}

	_, _ = b.WriteRune('\n')

	return b.String()
}

func Append(err error, errs ...error) *MultiError {
	if err, ok := err.(*MultiError); ok {
		err.errs = append(err.errs, errs...)

		return err
	}

	multiErr := &MultiError{
		errs: make([]error, 0, len(errs)+1),
	}

	if err != nil {
		multiErr.errs = append(multiErr.errs, err)
	}

	for _, err := range errs {
		if err != nil {
			multiErr.errs = append(multiErr.errs, err)
		}
	}

	return multiErr
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}

package main

import (
	"reflect"
	"slices"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Defragment(memory []byte, pointers []unsafe.Pointer) {
	pLen := len(pointers)
	nextToMove := 0

	for i := range memory {
		// memory defragmented
		if nextToMove == pLen {
			return
		}

		p := unsafe.Pointer(&memory[i])

		// skip if address in use
		if slices.Contains(pointers, p) {
			continue
		}

		// if pointer[d] is less than p, it means that pointer[d] is already in the defragmented area
		for nextToMove < pLen && uintptr(pointers[nextToMove]) < uintptr(p) {
			nextToMove++
		}

		// memory defragmented
		if nextToMove == pLen {
			return
		}

		memory[i] = *(*byte)(pointers[nextToMove])
		*(*byte)(pointers[nextToMove]) = 0
		pointers[nextToMove] = p

		nextToMove++
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

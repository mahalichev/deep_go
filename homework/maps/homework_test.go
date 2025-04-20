package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go
type node struct {
	key   int
	value int
	left  *node
	right *node
}

func newNode(key, value int) *node {
	return &node{
		key:   key,
		value: value,
	}
}

// true - new node, false - node was in bst
func (n *node) insert(key, value int) bool {
	if n.key == key {
		n.value = value

		return false
	}

	if key < n.key {
		if n.left != nil {
			return n.left.insert(key, value)
		}

		n.left = newNode(key, value)

		return true
	}

	if n.right != nil {
		return n.right.insert(key, value)
	}

	n.right = newNode(key, value)

	return true
}

// true - node was erased, false - there was no node with such a key in bst
func erase(n *node, key int) (*node, bool) {
	if n == nil {
		return nil, false
	}

	if key < n.key {
		newLeft, erased := erase(n.left, key)
		n.left = newLeft

		return n, erased
	}

	if key > n.key {
		newRight, erased := erase(n.right, key)
		n.right = newRight

		return n, erased
	}

	if n.left == nil {
		return n.right, true
	}

	if n.right == nil {
		return n.left, true
	}

	// nil value was checked earlier
	minN := minNode(n.right)

	n.key = minN.key
	n.value = minN.value

	newRight, erased := erase(n.right, key)
	n.right = newRight

	return n, erased
}

func minNode(n *node) *node {
	if n == nil {
		return nil
	}

	for n.left != nil {
		n = n.left
	}

	return n
}

func forEach(n *node, action func(int, int)) {
	if n == nil {
		return
	}

	forEach(n.left, action)
	action(n.key, n.value)
	forEach(n.right, action)
}

type OrderedMap struct {
	head *node
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

func (m *OrderedMap) Insert(key, value int) {
	if m.head == nil {
		m.head = newNode(key, value)
		m.size = 1

		return
	}

	if m.head.insert(key, value) {
		m.size++
	}
}

func (m *OrderedMap) Erase(key int) {
	head, erased := erase(m.head, key)

	if erased {
		m.head = head
		m.size--
	}
}

func (m *OrderedMap) Contains(key int) bool {
	for c := m.head; c != nil; {
		switch {
		case key < c.key:
			c = c.left
		case key > c.key:
			c = c.right
		default:
			return true
		}
	}

	return false
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	forEach(m.head, action)
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

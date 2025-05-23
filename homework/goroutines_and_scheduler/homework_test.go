package main

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

type wrappedTask struct {
	Task
	index    int
	priority int
}

type Heap []*wrappedTask

func (h Heap) Len() int {
	return len(h)
}

func (h Heap) Less(i, j int) bool {
	return h[i].priority > h[j].priority
}

func (h Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *Heap) Push(task any) {
	t := task.(*wrappedTask)
	t.index = len(*h)

	*h = append(*h, t)
}

func (h *Heap) Pop() any {
	n := len(*h)
	x := (*h)[n-1]
	*h = (*h)[:n-1]

	return x
}

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	h        *Heap
	idToTask map[int]*wrappedTask
}

func NewScheduler() Scheduler {
	h := make(Heap, 0)

	return Scheduler{
		h:        &h,
		idToTask: make(map[int]*wrappedTask),
	}
}

func (s *Scheduler) AddTask(task Task) {
	if _, exist := s.idToTask[task.Identifier]; exist {
		return
	}

	wrapped := &wrappedTask{
		Task:     task,
		priority: task.Priority,
	}

	s.idToTask[wrapped.Identifier] = wrapped
	heap.Push(s.h, wrapped)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	if wrapped, exist := s.idToTask[taskID]; exist {
		wrapped.priority = newPriority
		heap.Fix(s.h, wrapped.index)
	}
}

func (s *Scheduler) GetTask() Task {
	if len(*s.h) == 0 {
		return Task{}
	}

	wrapped := heap.Pop(s.h).(*wrappedTask)
	delete(s.idToTask, wrapped.Identifier)

	return wrapped.Task
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}

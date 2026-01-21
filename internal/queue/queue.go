package queue

import (
	"context"
)

type Action int

const (
	UnknownAction Action = 0
	InsertAction  Action = 1
	SaveAction    Action = 2
	DeleteAction  Action = 3
)

func (a Action) String() string {
	switch a {
	case InsertAction:
		return "InsertAction"
	case SaveAction:
		return "SaveAction"
	case DeleteAction:
		return "DeleteAction"

	}

	return "UnknownAction"
}

type SaveFunc func(index int64, body any) error
type DeleteFunc func(index int64) error
type ResultFunction func(action Action, index int64, err error)

func DefaultResultFunction(_ Action, _ int64, _ error) {}

type Queue struct {
	All   []*Handler
	Build *Build
}

func NewQueue(ctx context.Context, number uint, save, insert SaveFunc, delete DeleteFunc, threshold uint32, buildFunc func()) *Queue {
	if number < 1 {
		number = 1
	}

	out := &Queue{
		All:   make([]*Handler, number),
		Build: NewBuild(threshold, buildFunc),
	}

	for i := range number {
		out.All[i] = NewHandler(ctx, save, insert, delete, out.IncrementAndRun)
	}

	return out
}

func (q *Queue) Run(action Action, index int64, body any, resultFunc ResultFunction) {
	if resultFunc == nil {
		resultFunc = DefaultResultFunction
	}

	item := &Item{
		action:     action,
		index:      index,
		body:       body,
		ResultFunc: resultFunc,
	}

	number := int(index / 10 % int64(len(q.All)))
	q.All[number].Run(item)
}

func (q *Queue) IncrementAndRun() {
	q.Build.IncrementAndRun()
}

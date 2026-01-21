package toplist

import (
	"github.com/iostrovok/toplist/internal/queue"
)

type ResultFunction func(action Action, index int64, err error)

type Action string

const (
	unknownAction Action = "unknown"
	InsertAction  Action = "insert"
	SaveAction    Action = "save"
	DeleteAction  Action = "delete"
)

// ToQueue converts the public Action type into the internal queue.Action type used by the processor.
func (a Action) ToQueue() queue.Action {
	switch a {
	case SaveAction:
		return queue.SaveAction
	case InsertAction:
		return queue.InsertAction
	case DeleteAction:
		return queue.DeleteAction
	default:
		return queue.UnknownAction
	}
}

// FromQueue converts an internal queue.Action back into the public Action type.
func FromQueue(a queue.Action) Action {
	switch a {
	case queue.SaveAction:
		return SaveAction
	case queue.DeleteAction:
		return DeleteAction
	case queue.InsertAction:
		return InsertAction
	default:
		return unknownAction
	}
}

package toplist

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/iostrovok/toplist/internal"
	"github.com/iostrovok/toplist/internal/queue"
)

const (
	MinInfinity = math.MinInt64
	MaxInfinity = math.MaxInt64
)

// List represents the top-level structure that manages the skip list.
// It holds a queue for concurrent operations and a head pointer to the top-level node.
type List struct {
	Queue     *queue.Queue
	lastCount int            // Count of element according last rebuild operation
	head      unsafe.Pointer // Pointer to top-level head Node

	// for debug goals only
	debugMx   sync.Mutex
	debugMap  map[int64]queue.Action
	debugMode bool
}

// New initializes a new List with MinInfinity and MaxInfinity sentinels as head and tail.
// Also initializes the operation queue with 20 workers.
func New() *List {
	tail := internal.NewNode(MaxInfinity, unsafe.Pointer(nil), nil)
	head := internal.NewNode(MinInfinity, unsafe.Pointer(tail), nil)

	sl := &List{
		debugMap: map[int64]queue.Action{},
		head:     unsafe.Pointer(head),
	}

	sl.Queue = queue.NewQueue(context.Background(), 20, sl.BaseSave, sl.BaseInsert, sl.BaseDelete, 500, sl.Build)

	return sl
}

// Head returns the current top-level head pointer of the skip list.
func (tl *List) Head() unsafe.Pointer {
	return tl.head
}

// Build reconstructs or rebalances the skip list from the base layer up and swaps in the new head.
func (tl *List) Build() {
	head := tl.Head()
	baseNode := internal.GetBaseLavelHead(head)
	
	//top, count := internal.Build(head)
	top, count := internal.BuildSpeed(baseNode)

	tl.setLastCount(count)
	atomic.SwapPointer(&head, top)
}

// Clean physically removes all nodes marked for deletion from the base level of the skip list.
// It returns the total count of nodes that were removed.
func (tl *List) Clean() int {
	head := tl.Head()
	baseNode := internal.GetBaseLavelHead(head)
	return internal.CleanMarked(baseNode)
}

// LastCount returns the number of elements detected during the last rebuild (Build) operation.
func (tl *List) LastCount() int {
	return tl.lastCount
}

// setLastCount updates the internal counter for the number of elements in the list.
func (tl *List) setLastCount(count int) {
	tl.lastCount = count
}

// Run enqueues an operation for the list (such as Save or Delete) into the queue for concurrent handling.
// Converts internal queue actions back to public result function types.
func (tl *List) Run(action Action, index int64, body any, resultFunc ResultFunction) {
	rf := func(action queue.Action, index int64, err error) {
		tl.SaveDebugMap(index, action)
		if resultFunc != nil {
			resultFunc(FromQueue(action), index, err)
		}
	}

	tl.Queue.Run(action.ToQueue(), index, body, rf)
}

// Insert adds a new element to the list at the specified index and waits for the operation to complete.
func (tl *List) Insert(index int64, body any) error {
	ch := make(chan error, 1)
	resultFunc := func(action queue.Action, index int64, err error) {
		tl.SaveDebugMap(index, action)
		ch <- err
	}

	tl.Queue.Run(queue.InsertAction, index, body, resultFunc)
	return <-ch
}

// Save inserts or updates a node in the skip list at the given index and waits for completion.
// Uses the operation queue for concurrency-safety.
func (tl *List) Save(index int64, body any) error {
	ch := make(chan error, 1)
	resultFunc := func(action queue.Action, index int64, err error) {
		tl.SaveDebugMap(index, action)
		ch <- err
	}

	tl.Queue.Run(queue.SaveAction, index, body, resultFunc)
	return <-ch
}

// BaseInsert directly performs a non-idempotent insertion into
func (tl *List) BaseInsert(index int64, body any) error {
	return internal.Add(tl.Head, index, body, false, 10)
}

// BaseSave attempts to save an item into the skip list directly, without using the async queue.
// Returns an error if the operation fails after multiple retries.
func (tl *List) BaseSave(index int64, body any) error {
	return internal.Add(tl.Head, index, body, true, 10)
}

// BaseDelete deletes the node with the given index directly (not through the queue).
func (tl *List) BaseDelete(index int64) error {
	return internal.Delete(tl.Head, index)
}

// Find searches for a node by index, returning it and a flag indicating success.
func (tl *List) Find(index int64) (node *internal.Node, find bool) {
	return internal.Find(tl.Head(), index)
}

// Delete removes a node by index in the skip list and waits for the operation to complete using the queue.
func (tl *List) Delete(index int64) error {
	ch := make(chan error, 1)
	resultFunc := func(action queue.Action, index int64, err error) {
		tl.SaveDebugMap(index, action)
		ch <- err
	}

	tl.Queue.Run(queue.DeleteAction, index, nil, resultFunc)
	return <-ch
}

// PrintList prints a string representation of all levels of the skip list.
func (tl *List) PrintList() {
	list := tl.ToStringList()
	fmt.Println(strings.Join(list, "\n"))
}

// ToStringList returns the skip list as a slice of strings, one per level,
// where each string contains node values at that level.
func (tl *List) ToStringList() []string {
	levels := internal.ToStringList(tl.Head())
	out := make([]string, len(levels))
	for i, level := range levels {
		out[i] = strings.Join(level, " ")
	}

	return out
}

// ToString returns all nodes in the skip list, by level, as a slice of []*internal.Node per level.
func (tl *List) ToString() [][]*internal.Node {
	return internal.ToList(tl.Head())
}

// CheckBase verifies the base layer of the skip list and returns value and index slices for debug/consistency checks.
func (tl *List) CheckBase() ([]string, []int64) {
	return internal.CheckBase(tl.Head())
}

func (tl *List) CheckDuplicate() map[int64]int {
	return internal.CheckDuplicate(tl.Head())
}

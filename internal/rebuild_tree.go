package internal

import (
	"math"
	"sync/atomic"
	"unsafe"
)

const (
	MinInfinity = math.MinInt64
	MaxInfinity = math.MaxInt64
)

func Build(head unsafe.Pointer) (unsafe.Pointer, int) {
	topHead, count := BuildLevel(head)
	for {
		nextHead, _ := BuildLevel(topHead)
		if nextHead == nil {
			break
		}

		topHead = nextHead
	}

	return topHead, count
}

func BuildLevel(head unsafe.Pointer) (unsafe.Pointer, int) {
	preHead := loadUnmarkPointer(head)
	if preHead == nil {
		// false start
		return nil, 0
	}

	// new level
	headNode := NewNode(MinInfinity, nil, nil)
	atomic.StorePointer(&headNode.down, head)

	next := headNode
	counter := 0

	for preHead != nil && preHead.next != nil {
		if isMarked(preHead.next) {
			preHead = loadUnmarkPointer(preHead.next)
			continue
		}

		counter++

		if counter%10 == 0 {
			added := NewNode(preHead.Index, unsafe.Pointer(nil), nil)
			atomic.StorePointer(&added.down, unsafe.Pointer(preHead))
			atomic.StorePointer(&next.next, unsafe.Pointer(added))
			next = added
		}

		preHead = loadUnmarkPointer(preHead.next)
	}

	if counter < 10 {
		return nil, 0
	}

	atomic.StorePointer(&next.next, unsafe.Pointer(NewNode(MaxInfinity, unsafe.Pointer(nil), nil)))
	return unsafe.Pointer(headNode), counter
}

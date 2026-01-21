package internal

import (
	"sync/atomic"
	"unsafe"
)

func BuildSpeed(head unsafe.Pointer) (unsafe.Pointer, int) {
	preHead := loadUnmarkPointer(head)
	if preHead == nil {
		// false start
		return nil, 0
	}

	pingCh := make(chan struct{}, 10)
	headNodeResult := make(chan unsafe.Pointer, 3)

	go BuildLevelSpeedPing(head, 0, pingCh, headNodeResult)
	pingCh <- struct{}{}
	close(pingCh)
	headNode := <-headNodeResult

	return headNode, 0
}

func BuildLevelSpeedPing(head unsafe.Pointer, level int, pingCh chan struct{}, headNodeResult chan unsafe.Pointer) {
	level++
	preHead := loadUnmarkPointer(head)
	if preHead == nil {
		// false start
		return
	}

	// new level
	headNode := NewNode(MinInfinity, nil, nil)
	headNodePtr := unsafe.Pointer(headNode)
	atomic.StorePointer(&headNode.down, unsafe.Pointer(preHead))

	next := headNode
	counter := 0

	addedCount := 0
	nextPingCh := make(chan struct{}, 10)
	defer close(nextPingCh)

L:
	for {
		select {
		case _, ok := <-pingCh:
			if !ok {
				break L
			}

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

					addedCount++
					if addedCount == 1 {
						go BuildLevelSpeedPing(headNodePtr, level, nextPingCh, headNodeResult)
					}

					if addedCount%10 == 0 {
						nextPingCh <- struct{}{}
					}
				}

				preHead = loadUnmarkPointer(preHead.next)
			}
		}
	}

	nextPingCh <- struct{}{}

	if counter > 1 {
		atomic.StorePointer(&next.next, unsafe.Pointer(NewNode(MaxInfinity, unsafe.Pointer(nil), nil)))
	}

	if counter < 10 {
		headNodeResult <- head
	}
}

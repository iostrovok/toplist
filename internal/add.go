package internal

import (
	"sync/atomic"
	"unsafe"

	"github.com/iostrovok/toplist/terrors"
)

// Add attempts to insert a new node with the given index and value into the skip list.
// If replaceIfExists is true and a node at the index exists, its value is updated instead of inserting a new node.
// The operation retries up to maxRetry times in case of contention or races.
// Returns an error only if maxRetry is exhausted without success.
func Add(head TopFunc, index int64, body any, replaceIfExists bool, maxRetry int) (e error) {
	newNode := NewNode(index, nil, body)
	for maxRetry >= 0 {
		maxRetry--

		// Find predecessor node and the current node at the position of insertion.
		predHead, currNode, currPtr := findPredecessor(head(), index)
		if predHead == nil {
			continue
		}

		// If predHead is past the target or is in a "deleting" state, skip this round.
		if predHead.Index >= index || currPtr == nil || isMarked(currPtr) {
			continue
		}

		// If an unmarked node with the target Index already exists,
		// update its value if replaceIfExists is enabled.
		if currNode.Index < index || isMarked(currNode.next) {
			continue
		}

		if currNode.Index == index {
			if replaceIfExists {
				// Update the existing node's value.
				currNode.Value = body
			}

			return
		}

		// Link the new node to the successor.
		atomic.StorePointer(&newNode.next, currPtr)
		newPointer := unsafe.Pointer(newNode)

		// Attempt to physically insert the new node after predHead using CAS.
		if atomic.CompareAndSwapPointer(&predHead.next, currPtr, newPointer) {
			return
		}

		// If CAS fails, a race or concurrent modification occurred. Retry.
	}

	// Too many retries: indicate failure.
	return terrors.MaxRetryReached
}

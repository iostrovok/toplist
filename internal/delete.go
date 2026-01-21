package internal

import (
	"sync/atomic"
)

// Delete removes the node with the given index from the skip list.
// It searches for the node with the given index, and if found, marks it as deleted and physically unlinks it.
// If no such node is found, or it is already deleted, returns nil (no error).
func Delete(head TopFunc, index int64) error {
	for {
		// findPredecessor finds the predecessor and current node pointers at the given index.
		// predNode: node with index < target
		// currNode: node with index >= target (may be the one to delete)
		// currPtr: pointer to currNode
		predNode, currNode, currPtr := findPredecessor(head(), index)

		// If any pointer is invalid, or already marked, or predNode is past the index, try again.
		if predNode == nil || currPtr == nil || isMarked(currPtr) || predNode.Index >= index {
			continue
		}

		// If no matching node at the index, just return without error.
		if currNode == nil || currNode.Index > index {
			// not found
			return nil
		}

		// If found the node to delete, mark as deleted and try to physically remove it from the list using CAS.
		if currNode.Index == index {
			// Mark as logically deleted.
			marked := markPointer(currNode.next)
			if currNode.next != nil && !isMarked(currNode.next) {
				_ = atomic.CompareAndSwapPointer(&currNode.next, currNode.next, marked)
			}

			// Try to physically remove it from the list using CAS.
			if atomic.CompareAndSwapPointer(&predNode.next, currPtr, currNode.next) {
				return nil
			}
		}
	}
}

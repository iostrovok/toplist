package internal

import (
	"sync/atomic"
	"unsafe"
)

func CleanMarked(predPtr unsafe.Pointer) int {
	counter := 0

	var currNode *Node
	var currPtr unsafe.Pointer

	// Load and unmark the predPtr to get the node.
	predNode := loadUnmarkPointer(predPtr)
	if predNode == nil {
		// This should not happen; starting pointer should be less than target index.
		return counter
	}

	for {
		currPtr = predNode.next
		currNode = loadUnmarkPointer(predNode.next)
		if currNode == nil {
			// This should not happen; indicates a structural problem.
			return counter
		}

		// Physically remove all marked (logically deleted) nodes after predNode.
		for isMarked(currNode.next) {
			// Try to CAS out deleted node.
			if !atomic.CompareAndSwapPointer(&predNode.next, currPtr, currNode.next) {
				break
			}

			counter++

			// Advance currPtr and currNode (still removing marked nodes).
			currPtr = currNode.next
			currNode = loadUnmarkPointer(predNode.next)
		}

		// If currNode at this level is at or past index, or at list end
		if currNode.next == nil {
			return counter
		}

		// Advance to the next node at this level and continue.
		predPtr = predNode.next
		predNode = currNode
	}
}

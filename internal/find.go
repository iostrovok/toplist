package internal

import (
	"strings"
	"sync/atomic"
	"unsafe"
)

// Find searches for a node with the specified index starting from currPtrIn.
// Returns the node and true if found, or nil/false otherwise.
func Find(currPtrIn unsafe.Pointer, index int64) (node *Node, find bool) {
	for {
		// Traverse to find the predecessor node down the structure.
		preNode, node, _ := findPredecessor(currPtrIn, index)
		// If the candidate node is less than index or is logically deleted, continue searching.
		if node != nil && node.Index == index && !isMarked(node.next) {
			// If CAS fails, may be waiting for logically deleted node removal.
			return node, true
		}

		// If the predecessor is nil, or its index is not less than the target, or next is nil or marked, continue searching.
		if preNode == nil || preNode.Index >= index || preNode.next == nil || isMarked(preNode.next) {
			continue
		}

		// Get the next unmarked node (potential match).
		node = loadUnmarkPointer(preNode.next)
		if node == nil {
			// This should not happen; indicates a structural problem.
			return nil, false
		}

		// If the candidate node is less than index or is logically deleted, continue searching.
		if node.Index < index || isMarked(node.next) {
			// If CAS fails, may be waiting for logically deleted node removal.
			continue
		}

		// If found a node with a greater index, the search failed.
		if node.Index > index {
			return nil, false
		}

		// If the index matches, return the node.
		return node, true
	}
}

// findPredecessor travels down the skip list to find nodes surrounding the target index.
// predPtr: Node pointer to start from.
// Returns:
//   - predNode: Node with the biggest index < requested (or nil)
//   - currNode: Node matching the requested index (or nil)
//   - currPtr: Pointer to currNode (or nil)
//
// Only non-deleted (non-marked) elements are returned.
func findPredecessor(predPtr unsafe.Pointer, index int64) (predNode, currNode *Node, currPtr unsafe.Pointer) {
	// Load and unmark the predPtr to get the node.
	predNode = loadUnmarkPointer(predPtr)
	if predNode == nil || predNode.Index >= index {
		// This should not happen; starting pointer should be less than target index.
		return nil, nil, nil
	}

	for {
		currPtr = predNode.next
		currNode = loadUnmarkPointer(predNode.next)
		if currNode == nil {
			// This should not happen; indicates a structural problem.
			return nil, nil, nil
		}

		// Physically remove all marked (logically deleted) nodes after predNode.
		for isMarked(currNode.next) {
			// Try to CAS out deleted node.
			if !atomic.CompareAndSwapPointer(&predNode.next, currPtr, currNode.next) {
				break
			}

			// Advance currPtr and currNode (still removing marked nodes).
			currPtr = currNode.next
			currNode = loadUnmarkPointer(predNode.next)
		}

		// If currNode at this level is at or past index, or at list end
		if currNode.Index >= index || currNode.next == nil {
			// If there is no lower level, return current predNode and currNode.
			if predNode.down == nil {
				return predNode, currNode, currPtr
			}
			// Otherwise, descend one level down in the skip list.
			return findPredecessor(predNode.down, index)
		}

		// Advance to the next node at this level and continue.
		predPtr = predNode.next
		predNode = currNode
	}
}

// GetBaseLavelHead returns a pointer to the head node of the base (bottom) level of the skip list.
// It repeatedly traverses the "down" pointers from the provided head node until no further "down" node exists.
func GetBaseLavelHead(head unsafe.Pointer) unsafe.Pointer {
	for {
		next := loadUnmarkPointer(head).down
		if next == nil {
			return head
		}
		head = next
	}
}

// findPredecessor travels down the skip list to find nodes surrounding the target index.
// predPtr: Node pointer to start from.
// Returns:
//   - predNode: Node with the biggest index < requested (or nil)
//   - currNode: Node matching the requested index (or nil)
//   - currPtr: Pointer to currNode (or nil)
//
// Only non-deleted (non-marked) elements are returned.
func findPredecessorDebug(predPtr unsafe.Pointer, index int64) (predNode, currNode *Node, currPtr unsafe.Pointer, debug string) {
	debugOut := make([]string, 0)

	// Load and unmark the predPtr to get the node.
	predNode = loadUnmarkPointer(predPtr)
	if predNode == nil || predNode.Index >= index {
		// This should not happen; starting pointer should be less than target index.
		return nil, nil, nil, strings.Join(debugOut, "; ")
	}

	for {
		currPtr = predNode.next
		currNode = loadUnmarkPointer(predNode.next)
		if currNode == nil {
			// This should not happen; indicates a structural problem.
			return nil, nil, nil, strings.Join(debugOut, "; ")
		}

		// Physically remove all marked (logically deleted) nodes after predNode.
		for isMarked(currNode.next) {
			// Try to CAS out deleted node.
			if !atomic.CompareAndSwapPointer(&predNode.next, currPtr, currNode.next) {
				break
			}

			// Advance currPtr and currNode (still removing marked nodes).
			currPtr = currNode.next
			currNode = loadUnmarkPointer(predNode.next)
		}

		// If currNode at this level is at or past index, or at list end
		if currNode.Index >= index || currNode.next == nil {
			// If there is no lower level, return current predNode and currNode.
			if predNode.down == nil {
				return predNode, currNode, currPtr, strings.Join(debugOut, "; ")
			}
			// Otherwise, descend one level down in the skip list.
			return findPredecessorDebug(predNode.down, index)
		}

		// Advance to the next node at this level and continue.
		predPtr = predNode.next
		predNode = currNode
	}
}

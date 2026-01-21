package internal

import (
	"sync/atomic"
	"unsafe"
)

const (
	markBit = uintptr(1)
)

// --- Pointer Marking Utilities --

// markPointer sets the mark bit (lowest bit) of the provided unsafe pointer.
func markPointer(p unsafe.Pointer) unsafe.Pointer {
	if (uintptr(p) & markBit) != 0 {
		return unsafe.Pointer(uintptr(p) | markBit)
	}

	return p
}

// unmarkPointer clears the mark bit from the provided unsafe pointer to get the original address.
func unmarkPointer(p unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) &^ markBit)
}

// isMarked returns true if the provided unsafe pointer has the mark bit set.
func isMarked(p unsafe.Pointer) bool {
	return p != nil && (uintptr(p)&markBit) != 0
}

// loadUnmarkPointer unmarks the provided pointer and atomically loads the Node it points to.
func loadUnmarkPointer(p unsafe.Pointer) *Node {
	if p == nil {
		return nil
	}

	predPtrUn := unmarkPointer(p)
	in := atomic.LoadPointer(&predPtrUn)
	if in == nil {
		return nil
	}

	return (*Node)(in)
}

package toplist

import (
	"github.com/iostrovok/toplist/internal/queue"
)

// SetDebugMode enables or disables the debug mode for recording operations.
func (tl *List) SetDebugMode(mode bool) {
	tl.debugMode = mode
}

// SaveDebugMap records the last action performed on a specific index for debugging purposes.
func (tl *List) SaveDebugMap(index int64, action queue.Action) {
	if tl.debugMode {
		tl.debugMx.Lock()
		tl.debugMap[index] = action
		tl.debugMx.Unlock()
	}
}

// DebugMap returns the internal map containing recorded actions for debugging.
func (tl *List) DebugMap() map[int64]queue.Action {
	return tl.debugMap
}

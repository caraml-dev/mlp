package util

import "sync"

// ConcurrentSlice is a slice that can be appended concurrently
type ConcurrentSlice struct {
	sync.RWMutex
	items []interface{}
}

// Append safely appends items to a slice
func (cs *ConcurrentSlice) Append(item interface{}) {
	cs.Lock()
	defer cs.Unlock()

	cs.items = append(cs.items, item)
}

// GetItems gets the items from the list, should only be used when there are no operations working on it.
func (cs *ConcurrentSlice) GetItems() []interface{} {
	return cs.items
}

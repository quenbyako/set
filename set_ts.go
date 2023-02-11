package set

import (
	"sync"

	"golang.org/x/exp/maps"
)

// Set defines a thread safe set data structure.
type Set[T comparable] struct {
	set[T]
	sync.RWMutex // we name it because we don't want to expose it
}

var _ interface {
	rwLocker
	Interface[int]
} = (*Set[int])(nil)

// New creates and initialize a new Set. It's accept a variable number of
// arguments to populate the initial set. If nothing passed a Set with zero
// size is created.
func newTS[T comparable]() Interface[T] {
	return &Set[T]{set: set[T]{make(map[T]struct{})}}
}

type rwLocker interface {
	RLock()
	RUnlock()
}

// Add includes the specified items (one or more) to the set. The underlying
// Set s is modified. If passed nothing it silently returns.
func (s *Set[T]) Add(items ...T) {
	if len(items) == 0 {
		return
	}

	s.Lock()
	defer s.Unlock()

	for _, item := range items {
		s.m[item] = null{}
	}
}

// Remove deletes the specified items from the set.  The underlying Set s is
// modified. If passed nothing it silently returns.
func (s *Set[T]) Remove(items ...T) {
	if len(items) == 0 {
		return
	}

	s.Lock()
	defer s.Unlock()

	for _, item := range items {
		delete(s.m, item)
	}
}

// Pop  deletes and return an item from the set. The underlying Set s is
// modified. If set is empty, nil is returned.
func (s *Set[T]) Pop() (T, bool) {
	s.RLock()
	for item := range s.m {
		s.RUnlock()
		s.Lock()
		delete(s.m, item)
		s.Unlock()
		return item, true
	}
	s.RUnlock()
	var t T
	return t, false
}

// Has looks for the existence of items passed. It returns false if nothing is
// passed. For multiple items it returns true only if all of  the items exist.
func (s *Set[T]) Has(items ...T) bool {
	// assume checked for empty item, which not exist
	if len(items) == 0 {
		return false
	}

	s.RLock()
	defer s.RUnlock()

	has := true
	for _, item := range items {
		if _, has = s.m[item]; !has {
			break
		}
	}
	return has
}

// Size returns the number of items in a set.
func (s *Set[T]) Size() int {
	s.RLock()
	defer s.RUnlock()

	l := len(s.m)
	return l
}

// Clear removes all items from the set.
func (s *Set[T]) Clear() {
	s.Lock()
	defer s.Unlock()

	s.m = make(map[T]struct{})
}

// IsEqual test whether s and t are the same in size and have the same items.
func (s *Set[T]) IsEqual(t Interface[T]) bool {
	s.RLock()
	defer s.RUnlock()

	// Force locking only if given set is threadsafe.
	if conv, ok := t.(rwLocker); ok {
		conv.RLock()
		defer conv.RUnlock()
	}

	// return false if they are no the same size
	if sameSize := len(s.m) == t.Size(); !sameSize {
		return false
	}

	equal := true
	t.Each(func(item T) bool {
		_, equal = s.m[item]
		return equal // if false, Each() will end
	})

	return equal
}

// IsSubset tests whether t is a subset of s.
func (s *Set[T]) IsSubset(t Interface[T]) bool {
	s.RLock()
	defer s.RUnlock()

	return t.Each(func(item T) bool {
		_, ok := s.m[item]
		return ok
	})
}

// Each traverses the items in the Set, calling the provided function for each
// set member. Traversal will continue until all items in the Set have been
// visited, or if the closure returns false.
func (s *Set[T]) Each(f func(item T) bool) bool {
	s.RLock()
	defer s.RUnlock()

	return s.set.Each(f)
}

// List returns a slice of all items. There is also StringSlice() and
// IntSlice() methods for returning slices of type string or int.
func (s *Set[T]) List() []T {
	s.RLock()
	defer s.RUnlock()

	return maps.Keys(s.m)
}

// Copy returns a new Set with a copy of s.
func (s *Set[T]) Copy() Interface[T] {
	u := newTS[T]()
	for item := range s.m {
		u.Add(item)
	}
	return u
}

// Merge is like Union, however it modifies the current set it's applied on
// with the given t set.
func (s *Set[T]) Merge(t Interface[T]) {
	s.Lock()
	defer s.Unlock()

	t.Each(func(item T) bool {
		s.m[item] = null{}
		return true
	})

}

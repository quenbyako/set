// Package set provides both threadsafe and non-threadsafe implementations of
// a generic set data structure. In the threadsafe set, safety encompasses all
// operations on one set. Operations on multiple sets are consistent in that
// the elements of each set used was valid at exactly one point in time
// between the start and the end of the operation.
package set

// Set is describing a Set. Sets are an unordered, unique list of values.
type Set[T any] interface {
	Add(items ...T)
	Remove(items ...T)
	Pop() (T, bool)
	Has(items ...T) bool
	// Size returns the number of items in a set.
	Size() int
	// Clear removes all items from the set.
	Clear()
	// IsEmpty reports whether the Set is empty.
	IsEmpty() bool
	// IsEqual test whether s and t are the same in size and have the same items.
	IsEqual(s Set[T]) bool
	IsSubset(s Set[T]) bool
	IsSuperset(s Set[T]) bool
	Each(func(T) bool) bool
	String() string
	List() []T
	Copy() Set[T]
	Merge(s Set[T])
	Separate(s Set[T])
}

// helpful to not write everywhere struct{}{}
type null = struct{}

// New creates and initalizes a new Set interface. Its single parameter
// denotes the type of set to create. Either ThreadSafe or
// NonThreadSafe. The default is ThreadSafe.
func New[T comparable]() Set[T]      { return newTS[T]() }
func NewNonTS[T comparable]() Set[T] { return newNonTS[T]() }
func NewAny[T any]() Set[T]          { panic("unimplemented") }
func NewAnyNonTS[T any]() Set[T]     { panic("unimplemented") }

// Union is the merger of multiple sets. It returns a new set with all the
// elements present in all the sets that are passed.
//
// The dynamic type of the returned set is determined by the first passed set's
// implementation of the New() method.
func Union[T any](set1, set2 Set[T], sets ...Set[T]) Set[T] {
	u := set1.Copy()
	set2.Each(func(item T) bool {
		u.Add(item)
		return true
	})
	for _, set := range sets {
		set.Each(func(item T) bool {
			u.Add(item)
			return true
		})
	}

	return u
}

// Difference returns a new set which contains items which are in in the first
// set but not in the others. Unlike the Difference() method you can use this
// function separately with multiple sets.
func Difference[T any](set1, set2 Set[T], sets ...Set[T]) Set[T] {
	s := set1.Copy()
	s.Separate(set2)
	for _, set := range sets {
		s.Separate(set) // seperate is thread safe
	}
	return s
}

// Intersection returns a new set which contains items that only exist in all given sets.
func Intersection[T any](set1, set2 Set[T], sets ...Set[T]) Set[T] {
	all := Union(set1, set2, sets...)
	result := Union(set1, set2, sets...)

	all.Each(func(item T) bool {
		if !set1.Has(item) || !set2.Has(item) {
			result.Remove(item)
		}

		for _, set := range sets {
			if !set.Has(item) {
				result.Remove(item)
			}
		}
		return true
	})
	return result
}

// SymmetricDifference returns a new set which s is the difference of items which are in
// one of either, but not in both.
func SymmetricDifference[T any](s, t Set[T]) Set[T] {
	u := Difference(s, t)
	v := Difference(t, s)
	return Union(u, v)
}

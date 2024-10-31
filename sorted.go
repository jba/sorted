// For iter/sorted package proposal.

// Package sorted provides operations over iterators whose values are sorted.
package sorted

import (
	"cmp"
	"iter"
)

// Merge returns an iterator over all elements of s1 and s2, in the same order.
// Both s1 and s2 must be sorted.
func Merge[T cmp.Ordered](s1, s2 iter.Seq[T]) iter.Seq[T] {
	return MergeFunc(s1, s2, cmp.Compare[T])
}

// MergeFunc returns an iterator over all elements of s1 and s2, in the same order.
// Both s1 and s2 must be sorted according to cmp.
func MergeFunc[T any](s1, s2 iter.Seq[T], cmp func(T, T) int) iter.Seq[T] {
	return func(yield func(T) bool) {
		next, stop := iter.Pull(s2)
		defer stop()

		e2, ok := next()
		for e1 := range s1 {
			for ok && cmp(e2, e1) < 0 {
				if !yield(e2) {
					return
				}
				e2, ok = next()
			}
			if !yield(e1) {
				return
			}
		}
		for ; ok; e2, ok = next() {
			if !yield(e2) {
				return
			}
		}
	}
}

// Union returns an iterator over all elements of s1 and s2, in the same order.
// Each element appears only once in the result.
// Both s1 and s2 must be sorted.
func Union[T cmp.Ordered](s1, s2 iter.Seq[T]) iter.Seq[T] {
	return UnionFunc(s1, s2, cmp.Compare[T])
}

// UnionFunc returns an iterator over all elements of s1 and s2, in the same order.
// Each element appears only once in the result.
// Both s1 and s2 must be sorted according to cmp.
func UnionFunc[T any](s1, s2 iter.Seq[T], cmp func(T, T) int) iter.Seq[T] {
	return func(y func(T) bool) {
		yield := uniqueYielder(y, cmp)

		for x := range MergeFunc(s1, s2, cmp) {
			if !yield(x) {
				return
			}
		}
	}
}

// Intersect returns an iterator over all elements of s1 that are also in s2, in the same order.
// Each element of s1 appears only once in the result.
// Both s1 and s2 must be sorted.
func Intersect[T cmp.Ordered](s1, s2 iter.Seq[T]) iter.Seq[T] {
	return IntersectFunc(s1, s2, cmp.Compare[T])
}

// IntersectFunc returns an iterator over all elements of s1 that are also in s2, in the same order.
// Each element of s1 appears only once in the result.
// Both s1 and s2 must be sorted according to cmp.
func IntersectFunc[T any](s1, s2 iter.Seq[T], cmp func(T, T) int) iter.Seq[T] {
	return func(y func(T) bool) {
		next, stop := iter.Pull(s2)
		defer stop()

		yield := uniqueYielder(y, cmp)

		e2, ok := next()
		for e1 := range s1 {
			for ok && cmp(e2, e1) < 0 {
				e2, ok = next()
			}
			for ok && cmp(e2, e1) == 0 {
				if !yield(e1) {
					return
				}
				e2, ok = next()
			}
		}
	}
}

// Subtract returns an iterator over all elements of s1 that are not in s2, in the same order.
// Each element of s1 appears only once in the result.
// Both s1 and s2 must be sorted.
func Subtract[T cmp.Ordered](s1, s2 iter.Seq[T]) iter.Seq[T] {
	return SubtractFunc(s1, s2, cmp.Compare[T])
}

// SubtractFunc returns an iterator over all elements of s1 that are not in s2, in the same order.
// Each element of s1 appears only once in the result.
// Both s1 and s2 must be sorted according to cmp.
func SubtractFunc[T any](s1, s2 iter.Seq[T], cmp func(T, T) int) iter.Seq[T] {
	return func(y func(T) bool) {
		next, stop := iter.Pull(s2)
		defer stop()

		yield := uniqueYielder(y, cmp)

		e2, ok := next()
		for e1 := range s1 {
			for ok && cmp(e2, e1) < 0 {
				e2, ok = next()
			}
			if ok && cmp(e2, e1) == 0 {
				continue
			}
			if !yield(e1) {
				return
			}
		}
	}
}

func uniqueYielder[T any](yield func(T) bool, cmp func(T, T) int) func(T) bool {
	var (
		prev T
		have bool
	)
	return func(x T) bool {
		if have && cmp(x, prev) == 0 {
			return true
		}
		prev = x
		have = true
		return yield(x)
	}
}

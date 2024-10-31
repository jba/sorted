package sorted

import (
	"iter"
	"reflect"
	"slices"
	"testing"
)

func TestMerge(t *testing.T) {
	for ss1 := range sortedSlices(3, 3) {
		for ss2 := range sortedSlices(3, 3) {
			got := slices.Collect(Merge(slices.Values(ss1), slices.Values(ss2)))
			want := simpleMerge(ss1, ss2)
			if !slices.Equal(got, want) {
				t.Errorf("Merge(%v, %v) = %v, want %v", ss1, ss2, got, want)
			}
		}
	}
}

func simpleMerge(s1, s2 []int) []int {
	r := slices.Concat(s1, s2)
	slices.Sort(r)
	return r
}

func TestUnion(t *testing.T) {
	for ss1 := range sortedSlices(3, 3) {
		for ss2 := range sortedSlices(3, 3) {
			got := slices.Collect(Union(slices.Values(ss1), slices.Values(ss2)))
			want := simpleUnion(ss1, ss2)
			if !slices.Equal(got, want) {
				t.Errorf("Union(%v, %v) = %v, want %v", ss1, ss2, got, want)
			}
		}
	}
}

func simpleUnion(s1, s2 []int) []int {
	return slices.Compact(simpleMerge(s1, s2))
}

func TestIntersect(t *testing.T) {
	for ss1 := range sortedSlices(3, 3) {
		for ss2 := range sortedSlices(3, 3) {
			got := slices.Collect(Intersect(slices.Values(ss1), slices.Values(ss2)))
			want := simpleIntersect(ss1, ss2)
			if !slices.Equal(got, want) {
				t.Errorf("Intersect(%v, %v) = %v, want %v", ss1, ss2, got, want)
			}
		}
	}
}

func simpleIntersect(s1, s2 []int) []int {
	m1 := mapset(s1)
	m2 := mapset(s2)

	var r []int
	for x := range m1 {
		if m1[x] && m2[x] {
			r = append(r, x)
		}
	}
	slices.Sort(r)
	return r
}

func TestSubtract(t *testing.T) {
	for ss1 := range sortedSlices(3, 3) {
		for ss2 := range sortedSlices(3, 3) {
			got := slices.Collect(Subtract(slices.Values(ss1), slices.Values(ss2)))
			want := simpleSubtract(ss1, ss2)
			if !slices.Equal(got, want) {
				t.Errorf("Subtract(%v, %v) = %v, want %v", ss1, ss2, got, want)
			}
		}
	}
}

func simpleSubtract(s1, s2 []int) []int {
	m1 := mapset(s1)
	m2 := mapset(s2)

	var r []int
	for x := range m1 {
		if m1[x] && !m2[x] {
			r = append(r, x)
		}
	}
	slices.Sort(r)
	return r
}

func mapset(s []int) map[int]bool {
	m := map[int]bool{}
	for _, x := range s {
		m[x] = true
	}
	return m
}

// sortedSlices generates all slices whose values are in [0, n)
// and where each value appears at most max times.
func sortedSlices(n, max int) iter.Seq[[]int] {
	return func(yield func([]int) bool) {
		if n == 0 {
			yield([]int{})
			return
		}
		for ss := range sortedSlices(n-1, max) {
			if !yield(ss) {
				return
			}
			r := ss
			for range max {
				r = append(r[:len(r):cap(r)], n-1)
				if !yield(r) {
					return
				}
			}
		}
	}
}

func TestSortedSlices(t *testing.T) {
	got := slices.Collect(sortedSlices(2, 2))
	want := [][]int{
		{},
		{1},
		{1, 1},
		{0},
		{0, 1},
		{0, 1, 1},
		{0, 0},
		{0, 0, 1},
		{0, 0, 1, 1},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got\n%v\nwant\n%v", got, want)
	}
}

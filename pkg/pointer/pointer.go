// Package pointer provides utilities that assist in working with pointers.
package pointer

// To returns a pointer to the given value
func To[T any](v T) *T { return &v }

// From dereferences the pointer if it is not nil or returns d
func From[T any](p *T, d T) T {
	if p != nil {
		return *p
	}
	return d
}

// ToSlice returns a slice of pointers to the given values.
func ToSlice[T any](values []T) []*T {
	if len(values) == 0 {
		return nil
	}
	ret := make([]*T, 0, len(values))
	for _, v := range values {
		v := v
		ret = append(ret, &v)
	}
	return ret
}

// FromSlice returns a slice of values to the given pointers, dropping any nils.
func FromSlice[T any](values []*T) []T {
	if len(values) == 0 {
		return nil
	}
	ret := make([]T, 0, len(values))
	for _, v := range values {
		if v != nil {
			ret = append(ret, *v)
		}
	}
	return ret
}

// Equal reports if p is a pointer to a value equal to v
func Equal[T comparable](p *T, v T) bool {
	if p == nil {
		return false
	}
	return *p == v
}

// ValueEqual reports if value of pointer referenced by p is equal to value of pointer referenced by q
func ValueEqual[T comparable](p *T, q *T) bool {
	if p == nil || q == nil {
		return p == q
	}
	return *p == *q
}

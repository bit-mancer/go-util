// Code generated by slices-gen -- DO NOT EDIT

package slices

// UintSlice is a slice of uint's.
type UintSlice []uint

// Contains determines if the provided value is present in the slice.
// Runtime: O(n) -- appropriate for small-n slices.
func (s UintSlice) Contains(value uint) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}

	return false
}

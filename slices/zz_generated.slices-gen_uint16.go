// Code generated by slices-gen -- DO NOT EDIT

package slices

// Uint16Slice is a slice of uint16's.
type Uint16Slice []uint16

// Contains determines if the provided value is present in the slice.
// Runtime: O(n) -- appropriate for small-n slices.
func (s Uint16Slice) Contains(value uint16) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}

	return false
}

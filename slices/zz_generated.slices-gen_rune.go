// Code generated by slices-gen -- DO NOT EDIT

package slices

// RuneSlice is a slice of rune's.
type RuneSlice []rune

// Contains determines if the provided value is present in the slice.
// Runtime: O(n) -- appropriate for small-n slices.
func (s RuneSlice) Contains(value rune) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}

	return false
}
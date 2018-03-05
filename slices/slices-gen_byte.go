// Code generated by slices-gen -- DO NOT EDIT
// generated at 2018-03-04 20:13:49.614243745 -0800 PST m=+0.001045946

package util

// ByteSlice is a slice of byte's.
type ByteSlice []byte

// Contains determines if the provided value is present in the slice.
// Runtime: O(n) -- appropriate for small-n slices.
func (s ByteSlice) Contains(value byte) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}

	return false
}

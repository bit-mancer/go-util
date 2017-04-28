package util

// Int64InSlice determines if a value is present in the provided slice.
// Runtime: O(n) -- appropriate for small-n slices
func Int64InSlice(value int64, slice []int64) bool {

	for i, length := 0, len(slice); i < length; i++ {
		if slice[i] == value {
			return true
		}
	}

	return false
}

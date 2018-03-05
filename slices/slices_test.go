package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Int64InSlice", func() {
	It("accepts a nil slice", func() {
		Expect(func() { Int64InSlice(42, nil) }).NotTo(Panic())
		Expect(Int64InSlice(42, nil)).To(BeFalse())
	})

	It("returns false if the value is not in the slice", func() {
		var i int64 = 42
		v1 := []int64{1, 2, 4, 8, 16, 32, 64}

		Expect(Int64InSlice(i, nil)).To(BeFalse())
		Expect(Int64InSlice(i, []int64{})).To(BeFalse())
		Expect(Int64InSlice(i, v1)).To(BeFalse())
	})

	It("returns true if the value is in the slice", func() {
		var i int64 = 42
		v2 := []int64{7, 11, 13, 42}

		Expect(Int64InSlice(i, v2)).To(BeTrue())
	})
})

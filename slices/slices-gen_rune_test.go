// Code generated by slices-gen -- DO NOT EDIT

package slices

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RuneSlice", func() {
	// More of a static assertion, but eh:
	It("is a slice of rune's", func() {
		var s RuneSlice = []rune{1, 2, 3}
		s[0] = 0 // prevent unused error
	})

	_ = Describe(".Contains", func() {
		It("returns true if the provided value is in the slice", func() {
			s := RuneSlice{1, 2, 3, 4, 5}
			Expect(s.Contains(4)).To(BeTrue())
		})

		It("returns false if the provided value is not in the slice", func() {
			s := RuneSlice{1, 2, 3, 4, 5}
			Expect(s.Contains(42)).To(BeFalse())
		})
	})
})

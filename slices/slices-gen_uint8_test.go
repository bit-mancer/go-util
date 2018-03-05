// Code generated by slices-gen -- DO NOT EDIT
// generated at 2018-03-04 20:21:42.742651625 -0800 PST m=+0.001004428

package slices

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Uint8Slice", func() {
	// More of a static assertion, but eh:
	It("is a slice of uint8's", func() {
		var s Uint8Slice = []uint8{1, 2, 3}
		s[0] = 0 // prevent unused error
	})

	_ = Describe(".Contains", func() {
		It("returns true if the provided value is in the slice", func() {
			s := Uint8Slice{1, 2, 3, 4, 5}
			Expect(s.Contains(4)).To(BeTrue())
		})

		It("returns false if the provided value is not in the slice", func() {
			s := Uint8Slice{1, 2, 3, 4, 5}
			Expect(s.Contains(42)).To(BeFalse())
		})
	})
})

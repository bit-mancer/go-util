package config

import (
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type ifaceTestA struct{}

func (i ifaceTestA) String() string {
	return "A"
}

type ifaceTestB struct {
	IFace fmt.Stringer
}

func (i ifaceTestB) String() string {
	return "B"
}

type mixedExportFields struct {
	ExportedS   string
	unexportedI int
}

var _ = Describe("config", func() {

	// Exhaustive whitebox testing of isZeroValue because reflection.

	_ = Describe("isZeroValue", func() {

		// TODO stacktrace on test failure isn't capturing line numbers correctly with the helpers :(

		expectTrue := func(value interface{}) {
			Expect(isZeroValue(reflect.TypeOf(value), reflect.ValueOf(value))).To(BeTrue())
		}

		expectFalse := func(value interface{}) {
			Expect(isZeroValue(reflect.TypeOf(value), reflect.ValueOf(value))).To(BeFalse())
		}

		It("returns true if a provided value is the zero-value for its type", func() {

			// Mostly exhaustive, but doesn't bother testing /unsafe stuff

			expectTrue(nil)

			expectTrue(false)
			expectFalse(true)

			expectTrue(int(0))
			expectFalse(int(1))

			expectTrue(int8(0))
			expectFalse(int8(1))

			expectTrue(int16(0))
			expectFalse(int16(1))

			expectTrue(int32(0))
			expectFalse(int32(1))

			expectTrue(int64(0))
			expectFalse(int64(1))

			expectTrue(uint(0))
			expectFalse(uint(1))

			expectTrue(uint8(0))
			expectFalse(uint8(1))

			expectTrue(uint16(0))
			expectFalse(uint16(1))

			expectTrue(uint32(0))
			expectFalse(uint32(1))

			expectTrue(uint64(0))
			expectFalse(uint64(1))

			expectTrue(uintptr(0))

			expectTrue(float32(0))
			expectFalse(float32(1.0))

			expectTrue(float64(0))
			expectFalse(float64(1.0))

			expectTrue(complex64(0))
			expectFalse(complex64(1))

			expectTrue(complex128(0))
			expectFalse(complex128(1))

			expectTrue([3]int{})
			expectFalse([3]int{1, 2, 3})

			expectTrue(chan interface{}(nil))
			expectFalse(make(chan interface{}, 1))

			expectTrue((func() bool)(nil))
			expectFalse(func() bool { return false })

			expectTrue(fmt.Stringer(nil))
			// roundabout interface tests -- it'll hit the struct, then discover the embedded Stringer (testIFaceA) when iterating through the ifaceTestB fields
			expectTrue(fmt.Stringer(ifaceTestB{}))
			expectFalse(fmt.Stringer(ifaceTestB{ifaceTestA{}}))

			expectTrue(map[string]bool(nil))
			expectFalse(make(map[string]bool))

			expectTrue([]int(nil))
			expectFalse(&[]int{1, 2, 3})

			expectTrue("")
			expectFalse("test")

			expectTrue(struct{}{})
			expectTrue(struct{ i int }{0})
			expectTrue(struct{ i int }{1}) // unexported fields are skipped
			expectTrue(struct{ I int }{0})
			expectFalse(struct{ I int }{1})

			// Test unexported fields (it shouldn't panic on these)
			expectTrue(mixedExportFields{})
			expectFalse(mixedExportFields{ExportedS: "test", unexportedI: 1})
		})
	})
})

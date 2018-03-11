package util_test

import (
	"sync/atomic"

	. "github.com/bit-mancer/go-util/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WorkerPool", func() {

	var callCount *uint32
	var tasks chan interface{}
	var onTask func(interface{})
	var pool *WorkerPool

	BeforeEach(func() {
		count := uint32(0)
		callCount = &count

		tasks = make(chan interface{}, 1)

		onTask = func(interface{}) {
			atomic.AddUint32(callCount, 1)
		}

		pool = NewWorkerPool(tasks, onTask)
	})

	AfterEach(func(done Done) {
		pool.Abandon()
		pool.Wait()
		close(done)
	}, 3)

	Describe("NewWorkerPool", func() {
		It("creates a pool of initial size 0", func() {
			Expect(pool.Size()).To(Equal(0))
		})
	})

	Describe("Add", func() {
		It("Adds workers to the pool", func() {
			pool.Add(1)
			tasks <- 1

			Eventually(func() uint32 {
				return atomic.LoadUint32(callCount)
			}).Should(Equal(uint32(1)))
		})

		It("Panics if called on an abandoned pool", func() {
			pool.Add(1)
			pool.Abandon()
			Expect(func() { pool.Add(1) }).To(Panic())
		})
	})

	Describe("Remove", func() {
		It("Removes workers from the pool", func() {
			pool.Add(1)
			tasks <- 1

			Eventually(func() uint32 {
				return atomic.LoadUint32(callCount)
			}).Should(Equal(uint32(1)))

			Expect(pool.Remove(1)).To(BeNil())
			Expect(pool.Size()).To(Equal(0))
		})

		It("can be called throughout the lifetime of the pool", func() {
			pool.Add(3)
			tasks <- 1

			Eventually(func() uint32 {
				return atomic.LoadUint32(callCount)
			}).Should(Equal(uint32(1)))

			Expect(pool.Remove(2)).To(BeNil())

			pool.Add(4)
			Expect(pool.Remove(5)).To(BeNil())
		})

		It("Returns an error if the count exceeds the pool size", func() {
			pool.Add(2)
			Expect(pool.Remove(3)).NotTo(BeNil())
		})

		It("Panics if called on an abandoned pool", func() {
			pool.Add(1)
			pool.Abandon()
			Expect(func() { pool.Remove(1) }).To(Panic())
		})
	})

	Describe("Size", func() {
		It("changes based on Adds and Removes", func() {
			Expect(pool.Size()).To(Equal(0))
			pool.Add(3)
			Expect(pool.Size()).To(Equal(3))
			pool.Remove(2)
			Expect(pool.Size()).To(Equal(1))
		})

		It("can be called on an abandoned pool without panicking", func() {
			pool.Add(1)
			Expect(pool.Size()).To(Equal(1))
			pool.Abandon()
			Expect(pool.Size()).To(Equal(1))
		})
	})

	// TODO need more of an integration-level test to properly vet this
	Describe("Abandon", func() {
		It("", func(done Done) {
			pool.Add(1)
			tasks <- 1
			pool.Abandon()
			pool.Wait()

			close(done)
		})
	})

	// TODO need more of an integration-level test to properly vet this
	Describe("Wait", func() {
		It("blocks until all workers in the pool have terminated", func(done Done) {
			pool.Add(1)
			tasks <- 1
			close(tasks)
			pool.Wait()

			close(done)
		})
	})
})

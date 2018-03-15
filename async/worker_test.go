package async_test

import (
	"sync"
	"sync/atomic"

	. "github.com/bit-mancer/go-util/async"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Worker", func() {

	Describe("NewWorker", func() {
		It("requires a task channel and a handler func", func() {

			tasks := make(chan interface{}, 1)
			defer close(tasks)

			var callCount uint32 = 0

			onTask := func(interface{}) {
				atomic.AddUint32(&callCount, 1)
			}

			_, err := NewWorker(nil, nil, nil)
			Expect(err).NotTo(BeNil())

			_, err = NewWorker(tasks, nil, nil)
			Expect(err).NotTo(BeNil())

			_, err = NewWorker(nil, onTask, nil)
			Expect(err).NotTo(BeNil())

			_, err = NewWorker(tasks, onTask, nil)
			Expect(err).To(BeNil())
			tasks <- 1

			Eventually(func() uint32 {
				return atomic.LoadUint32(&callCount)
			}).Should(Equal(uint32(1)))
		})

		It("can optionally take a WaitGroup", func() {

			tasks := make(chan interface{}, 1)
			defer close(tasks)

			var callCount uint32 = 0

			onTask := func(interface{}) {
				atomic.AddUint32(&callCount, 1)
			}

			wg := &sync.WaitGroup{}

			_, err := NewWorker(tasks, onTask, wg)
			Expect(err).To(BeNil())
			tasks <- 1

			Eventually(func() uint32 {
				return atomic.LoadUint32(&callCount)
			}).Should(Equal(uint32(1)))
		})

		It("starts and returns a new Worker", func() {

			tasks := make(chan interface{}, 1)
			defer close(tasks)

			var callCount uint32 = 0

			onTask := func(interface{}) {
				atomic.AddUint32(&callCount, 1)
			}

			_, err := NewWorker(tasks, onTask, nil)
			Expect(err).To(BeNil())
			tasks <- 1

			Eventually(func() uint32 {
				return atomic.LoadUint32(&callCount)
			}).Should(Equal(uint32(1)))
		})
	})

	// TODO need more of an integration-level test to properly vet this
	It("can be waited upon to exit using a WaitGroup provided to NewWorker", func(done Done) {

		signal := make(chan struct{})
		tasks := make(chan interface{})

		onTask := func(interface{}) {
			signal <- struct{}{}
		}

		wg := sync.WaitGroup{}

		_, err := NewWorker(tasks, onTask, &wg)
		Expect(err).To(BeNil())
		tasks <- 1
		<-signal     // wait for worker to pick up the task
		close(tasks) // close the channel, signaling worker to exit

		wg.Wait()

		close(done)
	})

	Describe("Abandon", func() {
		// TODO need more of an integration-level test to properly vet this
		It("abandons any remaining tasks and causes the Worker to exit", func(done Done) {

			tasks := make(chan interface{})
			defer close(tasks)

			onTask := func(interface{}) {}

			wg := sync.WaitGroup{}

			worker, err := NewWorker(tasks, onTask, &wg)
			Expect(err).To(BeNil())
			worker.Abandon()
			wg.Wait()

			close(done)
		})
	})

	Describe("Wait", func() {
		// TODO need more of an integration-level test to properly vet this
		It("blocks until the Worker has exited", func(done Done) {

			signal := make(chan struct{})
			tasks := make(chan interface{})

			onTask := func(interface{}) {
				signal <- struct{}{}
			}

			w, err := NewWorker(tasks, onTask, nil) // no WaitGroup provided
			Expect(err).To(BeNil())
			tasks <- 1
			<-signal     // wait for worker to pick up the task
			close(tasks) // close the channel, signaling worker to exit

			w.Wait()

			close(done)
		})

		// TODO need more of an integration-level test to properly vet this
		It("also works with a WaitGroup provided to NewWorker", func(done Done) {

			signal := make(chan struct{})
			tasks := make(chan interface{})

			onTask := func(interface{}) {
				signal <- struct{}{}
			}

			wg := sync.WaitGroup{}

			_, err := NewWorker(tasks, onTask, &wg)
			Expect(err).To(BeNil())
			tasks <- 1
			<-signal     // wait for worker to pick up the task
			close(tasks) // close the channel, signaling worker to exit

			wg.Wait()

			close(done)
		})
	})
})

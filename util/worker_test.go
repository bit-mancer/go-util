package util_test

import (
	"sync"
	"sync/atomic"

	. "github.com/bit-mancer/go-util/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Worker", func() {
	It("is started by StartNewWorker", func() {

		tasks := make(chan interface{}, 1)
		defer func() { close(tasks) }()

		var callCount uint32 = 0

		onTask := func(interface{}) {
			atomic.AddUint32(&callCount, 1)
		}

		NewWorker(tasks, onTask, nil)
		tasks <- 1

		Eventually(func() uint32 {
			return atomic.LoadUint32(&callCount)
		}).Should(Equal(uint32(1)))
	})

	// TODO need more of an integration-level test to properly vet this
	It("can be waited upon to exit", func(done Done) {

		signal := make(chan struct{})
		tasks := make(chan interface{})

		onTask := func(interface{}) {
			signal <- struct{}{}
		}

		wg := sync.WaitGroup{}

		NewWorker(tasks, onTask, &wg)
		tasks <- 1
		<-signal     // wait for worker to pick up the task
		close(tasks) // close the channel, signaling worker to exit

		wg.Wait()

		close(done)
	})

	// TODO need more of an integration-level test to properly vet this
	It("can be abandoned", func(done Done) {

		tasks := make(chan interface{})
		defer func() { close(tasks) }()

		onTask := func(interface{}) {}

		wg := sync.WaitGroup{}

		worker := NewWorker(tasks, onTask, &wg)
		worker.Abandon()
		wg.Wait()

		close(done)
	})
})

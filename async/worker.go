package async

import (
	"fmt"
	"sync"

	"github.com/bit-mancer/go-util/util"
)

// Worker represents a goroutine that handles abstract, structured tasks. Workers can be pooled and managed via WorkerPool.
type Worker struct {
	_ util.NoCopy // trigger go vet on copy

	tasks      chan interface{}
	handleTask func(interface{})
	waitGroup  *sync.WaitGroup
	abandon    chan util.Signal
}

// NewWorker creates, starts, and returns a new Worker. The worker will accept items from the 'tasks' channel and run
// them by calling handleTask(). The worker will run until the 'tasks' channel is closed, or Abandon() is called.
//
// NewWorker will return an error if 'tasks' or 'handleTask' are nil.
// If 'waitGroup' is provided, Add/Done are called as the goroutine starts and stops; waitGroup can be nil if you
// don't need the start/stop information.
func NewWorker(tasks chan interface{}, handleTask func(interface{}), waitGroup *sync.WaitGroup) (*Worker, error) {

	if tasks == nil {
		return nil, fmt.Errorf("tasks channel cannot be nil")
	}

	if handleTask == nil {
		return nil, fmt.Errorf("handleTask func cannot be nil")
	}

	w := &Worker{
		tasks:      tasks,
		handleTask: handleTask,
		waitGroup:  waitGroup,
		abandon:    make(chan util.Signal)}

	start(w)
	return w, nil
}

// Design notes: having a signaling channel for quits and using a select, rather than doing a 'for range' on just the
// tasks channel, allows for the following:
//	- A portion of workers can be gracefully removed (via Abandon) in designs that use a single channel spread across
// 	  multiple workers, while the remaining workers continue to process the tasks in the channel.
//  - Closing the channel acts as a drain (the workers will run until they have consumed all the tasks), and the drain
// 	  can be interrupted by Abandon().
func start(w *Worker) {

	if w.waitGroup != nil {
		w.waitGroup.Add(1)
	}

	go func() {
		if w.waitGroup != nil {
			defer func() {
				w.waitGroup.Done()
			}()
		}

		for {
			select {
			case task, ok := <-w.tasks:
				if ok {
					w.handleTask(task)
				} else {
					return
				}

			case <-w.abandon:
				return
			}
		}
	}()
}

// Abandon instructs the worker goroutine to stop in the near future, possibly abandoning any remaining items in
// the worker task channel. Abandon is non-blocking and will immediately return, likely before the goroutine has
// stopped; use the WaitGroup passed to NewWorker to wait for the goroutine to actually stop.
//
// Note that an abandoned worker may complete another task before actually stopping.
//
// Abandon is not typically called to stop workers; instead, simply close the task channel (which acts as a
// drain -- no further tasks will be queued, any tasks left in the channel will be processed, then the worker(s)
// will exit).
func (w *Worker) Abandon() {
	go func() {
		w.abandon <- util.Signal{}
	}()
}

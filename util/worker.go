package util

import (
	"sync"
)

// Worker is implemented by goroutine-based workers.
type Worker interface {

	// Abandon instructs the worker goroutine to stop in the near future, possibly abandoning any remaining items in
	// the worker task channel. Abandon is non-blocking and will immediately return, likely before the goroutine has
	// stopped.
	//
	// Note that on average, the worker will complete one further task before actually stopping.
	// Abandon is not typically called to stop workers; instead, simply close the task channel (which acts as a
	// drain -- no further tasks will be queued, any tasks left in the channel will be processed, then the worker(s)
	// will exit).
	Abandon()
}

// worker represents a goroutine that handles abstract, structured tasks. Workers can be pooled and managed via
// WorkerPool.
type worker struct {
	_ NoCopy // trigger go vet on copy

	tasks     chan interface{}
	onTask    func(interface{})
	waitGroup *sync.WaitGroup
	abandon   chan Signal
}

// StartNewWorker creates, starts, and returns a new Worker.
// The worker will run until the 'tasks' channel is closed, or Abandon() is called.
func StartNewWorker(tasks chan interface{}, onTask func(interface{}), waitGroup *sync.WaitGroup) Worker {
	w := &worker{
		tasks:     tasks,
		onTask:    onTask,
		waitGroup: waitGroup,
		abandon:   make(chan Signal)}

	w.start()
	return w
}

/**
 * Design notes: having a signaling channel for quits and using a select, rather than doing a 'for range' on just the
 * tasks channel, allows for the following:
 *  - A portion of workers can be gracefully removed (via Abandon) in designs that use a single channel spread across
 *	  multiple workers, while the remaining workers continue to process the tasks in the channel.
 *  - Closing the channel acts as a drain (the workers will run until they have consumed all the tasks), and the drain
 *	  can be interrupted by Abandon().
 */
func (w *worker) start() {

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
					w.onTask(task)
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
// stopped.
//
// Note that on average, the worker will complete one further task before actually stopping.
// Abandon is not typically called to stop workers; instead, simply close the task channel (which acts as a
// drain -- no further tasks will be queued, any tasks left in the channel will be processed, then the worker(s)
// will exit).
func (w *worker) Abandon() {
	go func() {
		w.abandon <- Signal{}
	}()
}

package util

import (
	"fmt"
	"sync"
)

// WorkerPooler is implemented by a pool of workers that are specified via WorkerSpec.
type WorkerPooler interface {
	Worker

	Add(count int)
	Remove(count int) error
	Size() int
	Wait()
}

type workerPool struct {
	// Worker spec:
	tasks     chan interface{}
	onTask    func(interface{})
	waitGroup *sync.WaitGroup

	// Mutex covers everything below:
	mutex sync.Mutex

	workers     []Worker
	isAbandoned bool
}

// NewWorkerPool returns a WorkerPooler whose workers receive work from the tasks channel, and perform the work by
// calling onTask() on the received work. The pool is initially empty.
// THREAD-SAFETY: the WorkerPooler is thread-safe.
func NewWorkerPool(tasks chan interface{}, onTask func(interface{})) WorkerPooler {
	return &workerPool{
		tasks:     tasks,
		onTask:    onTask,
		waitGroup: &sync.WaitGroup{},
		mutex:     sync.Mutex{},
		workers:   make([]Worker, 0)}
}

// Add creates, starts, and adds to the pool a number of workers equal to count.
// Attempting to add to an abandoned pool will result in a panic.
func (p *workerPool) Add(count int) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.isAbandoned {
		panic("Trying to add workers after pool has been abandoned!")
	}

	newWorkers := make([]Worker, count)
	for i := range newWorkers {
		newWorkers[i] = StartNewWorker(p.tasks, p.onTask, p.waitGroup)
	}

	p.workers = append(p.workers, newWorkers...)
}

// Remove abandons and removes from the pool a number of workers equal to count; abandoned workers will continue to
// process any current task, and may continue to run for a short period of time.
// Remove returns an error if count exceeds the size of the pool (no workers will be removed in this case).
// Attempting to remove from an abandoned pool will result in a panic.
func (p *workerPool) Remove(count int) error {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.isAbandoned {
		panic("Trying to remove workers after pool has been abandoned!")
	}

	length := len(p.workers)

	if count > length {
		return fmt.Errorf("tried to remove more workers (%d) than were available in pool (%d)", count, length)
	}

	firstIndexToRemove := length - count
	for i := firstIndexToRemove; i < length; i++ {
		p.workers[i].Abandon()
		p.workers[i] = nil
	}

	p.workers = p.workers[:firstIndexToRemove]

	return nil
}

// Size returns the number of workers in the pool.
func (p *workerPool) Size() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return len(p.workers)
}

// Abandon instructs all workers in the pool to stop in the near future, possibly abandoning any remaining items in the
// worker task channel. Abandon is non-blocking and will immediately return, likely before the workers have stopped.
//
// Note that on average, each worker will complete one further task before actually stopping.
// Abandon is not typically called to stop workers; instead, simply close the task channel (which acts as a drain --
// no further tasks will be queued, and any tasks left in the channel will be processed, then the worker(s) will exit).
func (p *workerPool) Abandon() {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.isAbandoned {
		return
	}
	p.isAbandoned = true

	for _, w := range p.workers {
		w.Abandon()
	}
}

// Wait is a blocking call that waits for all workers in the pool to stop.
// IMPORTANT: You must have closed the task channel and/or called Abandon() prior to calling Wait, otherwise a
// deadlock will occur.
func (p *workerPool) Wait() {
	// Don't need the mutex
	p.waitGroup.Wait()
}

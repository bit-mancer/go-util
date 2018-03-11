package async

import (
	"fmt"
	"sync"
)

// WorkerPool is a pool of goroutine-based Workers.
type WorkerPool struct {
	// Worker spec:
	tasks     chan interface{}
	onTask    func(interface{})
	waitGroup *sync.WaitGroup

	// Mutex covers everything below:
	mutex sync.Mutex // struct will be no-copy due to the mutex

	workers     []*Worker
	isAbandoned bool
}

// NewWorkerPool returns a WorkerPool whose workers receive work from the provided tasks channel, and perform the work
// by calling onTask() on the received work. The pool is initially empty.
// THREAD-SAFETY: the WorkerPool is thread-safe.
func NewWorkerPool(tasks chan interface{}, onTask func(interface{})) *WorkerPool {
	return &WorkerPool{
		tasks:     tasks,
		onTask:    onTask,
		waitGroup: &sync.WaitGroup{},
		mutex:     sync.Mutex{},
		workers:   make([]*Worker, 0)}
}

// Add creates, starts, and adds to the pool a number of workers equal to count.
// Attempting to add to an abandoned pool will result in a panic.
func (p *WorkerPool) Add(count int) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.isAbandoned {
		panic("Trying to add workers after pool has been abandoned!")
	}

	newWorkers := make([]*Worker, count)
	for i := range newWorkers {
		newWorkers[i] = NewWorker(p.tasks, p.onTask, p.waitGroup)
	}

	p.workers = append(p.workers, newWorkers...)
}

// Remove abandons and removes from the pool a number of workers equal to count. Abandoned workers will complete any
// current task they have, and may pick up another task before actually stopping (on average, half the abandoned
// workers will run another task before stopping).
// Remove returns an error if count exceeds the size of the pool (no workers will be removed in this case).
// Attempting to remove from an abandoned pool will result in a panic.
func (p *WorkerPool) Remove(count int) error {

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
func (p *WorkerPool) Size() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return len(p.workers)
}

// Abandon instructs all workers in the pool to stop in the near future, possibly abandoning any remaining items in the
// worker task channel. Abandon is non-blocking and will immediately return, likely before the workers have stopped;
// use Wait() to wait for all the workers to actually stop.
//
// Note that on average, half the workers will complete one further task before actually stopping.
// Abandon is not typically called to stop workers; instead, simply close the task channel (which acts as a drain --
// no further tasks will be queued, and any tasks left in the channel will be processed, then the workers will exit).
func (p *WorkerPool) Abandon() {

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
func (p *WorkerPool) Wait() {
	// Don't need the mutex
	p.waitGroup.Wait()
}

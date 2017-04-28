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

// NewWorkerPool returns a new WorkerPool whose workers are specified by the provided WorkerSpec. The pool is initially
// empty.
// THREAD-SAFETY: the WorkerPool instance is thread-safe.
func NewWorkerPool(tasks chan interface{}, onTask func(interface{})) WorkerPooler {
	return &workerPool{
		tasks:     tasks,
		onTask:    onTask,
		waitGroup: &sync.WaitGroup{},
		mutex:     sync.Mutex{},
		workers:   make([]Worker, 0)}
}

// Add creates, starts, and adds to the pool a number of workers equal to count.
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

// Remove stops and removes the pool a number of workers equal to count.
// Remove returns an error if count exceeds the size of the pool (no workers will be removed in this case).
func (p *workerPool) Remove(count int) error {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.isAbandoned {
		panic("Trying to remove workers after pool has been abandoned!")
	}

	if count > len(p.workers) {
		return fmt.Errorf("Tried to remove more workers (%d) than were available in pool (%d).", count, len(p.workers))
	}

	for i := len(p.workers) - count; i < len(p.workers); i++ {
		p.workers[i].Abandon()
		p.workers[i] = nil
	}

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
// Pass a sync.WaitGroup as part of the WorkerSpec if you need to wait for the workers to stop.
//
// Note that on average, each worker will complete one further task before actually stopping.
// Abandon is not typically called to stop workers; instead, simply close the task channel (which acts as a drain --
// no further tasks will be queued, any any tasks left in the channel will be processed, then the worker(s) will exit).
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

// Wait is a blocking call that waits for all workers in the pool to stop. You must have closed the task channel and/or
// called Abandon() prior to calling Wait, otherwise a deadlock will occur.
func (p *workerPool) Wait() {
	// Don't need the mutex
	p.waitGroup.Wait()
}

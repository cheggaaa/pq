package queue

import (
	"container/heap"
	"fmt"
	"sync"
	"sync/atomic"
	"runtime"
)

var ErrQueueNotStarted = fmt.Errorf("Queue not started or closed")
var ErrQueueAlreayStarted = fmt.Errorf("Queue already started")

type Queue struct {
	numWorkers int
	pq priorityQueue
	work chan *item
	m sync.Mutex
	wg sync.WaitGroup
	working bool
	taskRunning int32
}

func (q *Queue) Start(numWorkers int) (err error) {
	q.m.Lock()
	defer q.m.Unlock()
	if q.working {
		return ErrQueueAlreayStarted
	}
	q.numWorkers = numWorkers
	q.pq = make(priorityQueue, 0)
	q.work = make(chan *item)
	if q.numWorkers <= 0 {
		q.numWorkers = runtime.NumCPU()
	}
	q.runWorkers()
	go q.dispatcher()
	return
}

func (q *Queue) AddFunc(f func(), priority int) (err error) {
	task := &funcTask{
		f : f,
		p : priority,
	}
	return q.AddTask(task)
}

func (q *Queue) AddFuncAndWait(f func(), priority int) (err error) {
	task := &funcTask{
		f : f,
		p : priority,
	}
	return q.AddTaskAndWait(task)
}

func (q *Queue) AddTask(task Task) (err error) {
	it := &item{task:task}
	return q.addItem(it)
}

func (q *Queue) AddTaskAndWait(task Task) (err error) {
	// add
	it := &item{task:task, done:make(chan bool)}
	if err = q.addItem(it); err != nil {
		return
	}
	// and wait
	<- it.done
	return
}

func (q *Queue) Len() int {
	return len(q.pq)
}

func (q *Queue) TaskRunning() int {
	return int(atomic.LoadInt32(&q.taskRunning))
}

func (q *Queue) addItem(it *item) (err error) {
	q.m.Lock()
	defer q.m.Unlock()
	if ! q.working {
		return ErrQueueNotStarted
	}
	heap.Push(&q.pq, it)
	if q.pq.Len() == 1 {
		go q.dispatcher()
	}
	return
}

func (q *Queue) runWorkers() {
	for i := 0; i < q.numWorkers; i++ {
		go q.worker()
	}
	q.working = true
}

func (q *Queue) dispatcher() {
	for {
		q.m.Lock()
		if q.pq.Len() == 0 || ! q.working {
			q.m.Unlock()
			break
		}
		it := heap.Pop(&q.pq)
		q.m.Unlock()
		q.work <- it.(*item)
	}
}

func (q *Queue) worker() {
	q.wg.Add(1)
	defer q.wg.Done()
	for it := range q.work {
		atomic.AddInt32(&q.taskRunning, 1)
		it.task.Run()
		atomic.AddInt32(&q.taskRunning, -1)
		if it.done != nil {
			it.done <- true
		}
	}
}

func (q *Queue) Stop() {
	q.m.Lock()
	defer q.m.Unlock()
	close(q.work)
	q.working = false
	q.wg.Wait()
}
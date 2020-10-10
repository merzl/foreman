package subscriber

import (
	"context"
)

type task interface {
	do()
}

type dispatcherQueue chan chan task
type workerQueue chan task

type worker struct {
	ctx             context.Context
	dispatcherQueue dispatcherQueue
	myTasks         workerQueue
}

func (w *worker) workerQueue() workerQueue {
	return w.myTasks
}

func newWorker(ctx context.Context, dispatcherQueue dispatcherQueue, myTasks workerQueue) worker {
	return worker{
		ctx: ctx,
		dispatcherQueue: dispatcherQueue,
		myTasks: myTasks,
	}
}

func(w *worker) start() {
	go func() {
		for {
			//tell dispatcher that I'm ready to work
			w.dispatcherQueue <- w.myTasks
			select {
			case <-w.ctx.Done():
				return
			case task, open := <- w.myTasks:
				if !open {
					return
				}
				task.do()
			}
		}
	}()
}

func newDispatcher(workersCount uint) *dispatcher {
	return &dispatcher{
		workersCount: workersCount,
		workersQueues: make(dispatcherQueue, workersCount),
		workersWorkplaces: make([]workerQueue, workersCount),
	}
}

type dispatcher struct {
	workersCount      uint
	workersQueues      dispatcherQueue
	workersWorkplaces []workerQueue
}

func (d *dispatcher) busyWorkers() int {
	return len(d.workersQueues)
}

func (d *dispatcher) start(ctx context.Context) {
	go func() {
		<- ctx.Done()
		d.workersQueues = nil
		for _, c := range d.workersWorkplaces {
			close(c)
		}
	}()

	var i uint
	for i < d.workersCount {
		d.workersWorkplaces[i] = make(workerQueue)
		worker := newWorker(ctx, d.workersQueues, d.workersWorkplaces[i])
		worker.start()
		i++
	}
}

func (d dispatcher) schedule(task task) {
	worker := <- d.workersQueues
	worker <- task
}

type dumpTask struct {

}
func (t *dumpTask) do() {

}
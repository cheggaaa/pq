package pq

import "sync/atomic"

// is a wrapper for Task
type item struct {
	task  Task
	index int
	ctrl  *itemCtrl
}

type itemCtrl struct {
	count  int32
	done   func(err error)
}

func (i *item) can() bool {
	if i.ctrl != nil {
		return atomic.LoadInt32(&i.ctrl.count) > 0
	}
	return true
}

func (i *item) done(err error) {
	if i.ctrl != nil {
		if atomic.AddInt32(&i.ctrl.count, -1) <= 0 || err != nil {
			if i.ctrl.done != nil {
				i.ctrl.done(err)
				i.ctrl.done = nil
			}
		}
		if err != nil {
			atomic.StoreInt32(&i.ctrl.count, -1)
		}
	}
}

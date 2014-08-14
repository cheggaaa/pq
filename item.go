package pq

import (
	"sync"
)

// is a wrapper for Task
type item struct {
	task  Task
	index int
	ctrl  *itemCtrl
}

type itemCtrl struct {
	count  int32
	done   chan error
	m *sync.Mutex
}

func (i *item) can() bool {
	if i.ctrl != nil {
		i.ctrl.m.Lock()
		is := i.ctrl.count > 0
		i.ctrl.m.Unlock()
		return is
	}
	return true
}

func (i *item) done(err error) {
	if i.ctrl != nil {
		i.ctrl.m.Lock()	
		i.ctrl.count--
		if i.ctrl.count <= 0 || err != nil {
			if i.ctrl.done != nil {
				i.ctrl.done <- err
				i.ctrl.done = nil
			}
		}
		if err != nil {
			i.ctrl.count = -1
		}
		i.ctrl.m.Unlock()
	}
}

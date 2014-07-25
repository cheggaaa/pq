package pq

// task interface
type Task interface {
	Priority() int
	Run() (err error)
}

// default task type
type funcTask struct {
	f func() error
	p int
}

func (ft *funcTask) Priority() int {
	return ft.p
}

func (ft *funcTask) Run() (err error) {
	return ft.f()
}

package queue

// task interface
type Task interface {
	Priority() int
	Run()
}

// default task type
type funcTask struct {
	f func()
	p int
}

func (ft *funcTask) Priority() int {
	return ft.p
}

func (ft *funcTask) Run() {
	ft.f()
}

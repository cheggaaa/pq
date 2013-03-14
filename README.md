## Simple priority queue
==

1. Create queue   
2. Start N workers   
3. Execute your tasks over workers (sync or async)   
4. Make high priority for important tasks   

### Installation:
```
go get -u github.com/cheggaaa/pq
```

### Example:
```Go
package main

import (
  	"fmt"
	"github.com/cheggaaa/pq"
	"sync"
	"time"
)

type HardWork struct {
	name string
	priority int
	duration int
}

// implement pq.Task
func (w *HardWork) Priority() int {
	return w.priority
}

func (w *HardWork) Run() {
	fmt.Println("Start:", w.name)
	time.Sleep(time.Duration(w.duration) * time.Second)
	fmt.Println("Done:", w.name)
}

var ToDo = []*HardWork{
	{
		name:     "Convert picture 1",
		priority: 5,
		duration: 1,
	},
	{
		name:     "Convert picture 2",
		priority: 5,
		duration: 1,
	},
	{
		name:     "Convert picture 3",
		priority: 5,
		duration: 1,
	},
	{
		name:     "Convert picture 4",
		priority: 5,
		duration: 1,
	},
	{
		name:     "Sing a song",
		priority: 100, // very important :-)
		duration: 10,
	},
	{
		name:     "Convert video 1",
		priority: 50,
		duration: 5,
	},
	{
		name:     "Convert video 2",
		priority: 55,
		duration: 5,
	},
}

func main() {
	wg := sync.WaitGroup{}

	// create queue
	q := &pq.Queue{}
	// start two workers
	q.Start(2)

	// for control
	wg.Add(len(ToDo))

	// add tasks
	for _, task := range ToDo {
		go func(task *HardWork) {
			fmt.Println("Add task:", task.name)
			q.AddTaskAndWait(task)
			wg.Done()
		}(task)
	}
	wg.Wait()
	q.Stop()
}

```



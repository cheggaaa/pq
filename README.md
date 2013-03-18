## Simple priority queue
==

1. Create queue   
2. Start N workers   
3. Execute your tasks over workers (sync or async) Create groups of tasks (if you want)
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
	"time"
)

type HardWork struct {
	name     string
	priority int
	duration int
}

// implement pq.Task
func (w *HardWork) Priority() int {
	return w.priority
}

func (w *HardWork) Run() {
	fmt.Printf("Start: %s (%d)\n", w.name, w.priority)
	time.Sleep(time.Duration(w.duration) * time.Second)
	fmt.Println("Done:", w.name)
}

var ToDo = []pq.Task{
	&HardWork{
		name:     "Convert picture 1",
		priority: 5,
		duration: 1,
	},
	&HardWork{
		name:     "Convert picture 2",
		priority: 5,
		duration: 1,
	},
	&HardWork{
		name:     "Convert picture 3",
		priority: 5,
		duration: 1,
	},
	&HardWork{
		name:     "Convert picture 4",
		priority: 5,
		duration: 1,
	},
	&HardWork{
		name:     "Sing a song",
		priority: 100, // very important :-)
		duration: 10,
	},
	&HardWork{
		name:     "Convert video 1",
		priority: 50,
		duration: 5,
	},
	&HardWork{
		name:     "Convert video 2",
		priority: 55,
		duration: 5,
	},
}

func main() {
	// create queue
	q := &pq.Queue{}
	// start two workers
	q.Start(2)

	// add and wait ToDo
	fmt.Println("Add tasks...")
	q.WaitGroup(ToDo)
	fmt.Println("ToDo done!")

	// just function
	q.WaitFunc(func() {
		fmt.Println("Just do this")
	}, 0)

	// add one more task
	fmt.Println(".. and last task")
	q.WaitTask(&HardWork{
		name:     "Last mega task",
		duration: 1,
	})
	fmt.Println("Now all tasks is done, stop queue")
	q.Stop()
}

```



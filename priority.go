package pq

// A PriorityQueue implements heap.Interface and holds Items.
type priorityQueue []*item

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].task.Priority() > pq[j].task.Priority()
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x interface{}) {
	item := x.(*item)
	q := *pq
	q = append(q, item)
	*pq = q
}

func (pq *priorityQueue) Pop() interface{} {
	a := *pq
	n := len(a)
	item := a[n-1]
	*pq = a[0 : n-1]
	return item
}

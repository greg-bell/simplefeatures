package rtree

import "container/heap"

// NearestSearch iterates over the records in the RTree in order of distance
// from the input box (shortest distanace first using the Euclidean metric).
// The callback is called for every element iterated over. If an error is
// returned from the callback, then iteration stops immediately. Any error
// returned from the callback is returned by NearestSearch, except for the case
// where the special Stop sentinal error is returned (in which case nil will be
// returned from NearestSearch).
func (t *RTree) NearestSearch(box Box, callback func(recordID int) error) error {
	if t.root == nil {
		return nil
	}

	queue := entriesQueue{origin: box}
	equeueNode := func(n *node) {
		for i := 0; i < n.numEntries; i++ {
			heap.Push(&queue, &n.entries[i])
		}
	}

	equeueNode(t.root)
	for len(queue.entries) > 0 {
		nearest := heap.Pop(&queue).(*entry)
		if nearest.child == nil {
			if err := callback(nearest.recordID); err != nil {
				if err == Stop {
					return nil
				}
				return err
			}
		} else {
			equeueNode(nearest.child)
		}
	}
	return nil
}

type entriesQueue struct {
	entries []*entry
	origin  Box
}

func (q *entriesQueue) Len() int {
	return len(q.entries)
}

func (q *entriesQueue) Less(i int, j int) bool {
	e1 := q.entries[i]
	e2 := q.entries[j]
	return squaredEuclideanDistance(e1.box, q.origin) < squaredEuclideanDistance(e2.box, q.origin)
}

func (q *entriesQueue) Swap(i int, j int) {
	q.entries[i], q.entries[j] = q.entries[j], q.entries[i]
}

func (q *entriesQueue) Push(x interface{}) {
	q.entries = append(q.entries, x.(*entry))
}

func (q *entriesQueue) Pop() interface{} {
	e := q.entries[len(q.entries)-1]
	q.entries = q.entries[:len(q.entries)-1]
	return e
}

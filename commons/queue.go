package commons

type Queue[T comparable] struct {
	data []T
}

func (q *Queue[T]) Enqueue(v T) {
	q.data = append(q.data, v)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	var zero T
	if len(q.data) == 0 {
		return zero, false
	}
	v := q.data[0]
	q.data = q.data[1:]
	return v, true
}

// check if the queue contains the given value
func (q *Queue[T]) Contains(v T) bool {
	for _, item := range q.data {
		if item == v {
			return true
		}
	}
	return false
}

func (q *Queue[T]) Len() int { return len(q.data) }

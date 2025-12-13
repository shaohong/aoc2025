package queue

type Queue[T any] struct {
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

func (q *Queue[T]) Len() int { return len(q.data) }

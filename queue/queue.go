package queue

import (
	"fmt"
	"jobQ/job"
	"jobQ/list"
	"sync"
	"time"
)

func New() *Queue {
	return &Queue{
		list:  list.New(),
		value: sync.Map{},
	}
}

type jobNP struct {
	job.Job
	np *list.Node
}

type Queue struct {
	list  list.List
	value sync.Map
}

func (q *Queue) Enqueue(j job.Job) error {
	if _, ok := q.value.Load(j.Id); ok {
		return ErrIdCollision
	}

	var np *list.Node
	if j.Status == job.Queued {
		np = q.list.Push(j.Id)
	}
	q.value.Store(j.Id, jobNP{
		Job: j,
		np:  np,
	})
	return nil
}

func (q *Queue) Dequeue() (job.Job, error) {
	id, err := q.list.Pop()
	if err != nil {
		return job.Job{}, err
	}

	jAny, ok := q.value.Load(id)
	if !ok {
		return job.Job{}, fmt.Errorf("%w job with id=%d", ErrNotFound, id)
	}

	j := jAny.(jobNP)
	j.Status = job.InProgress
	j.np = nil
	q.value.Store(j.Id, j)

	return j.Job, nil
}

func (q *Queue) Conclude(id uint64) error {
	jAny, ok := q.value.Load(id)
	if !ok {
		return fmt.Errorf("%w job with id=%d", ErrNotFound, id)
	}
	j := jAny.(jobNP)

    if j.Status == job.Concluded {
        return nil
    }

    if j.Status == job.Canceled {
        return fmt.Errorf("%w job with id=%d was already canceled", ErrFinalState, id)
    }

	if j.np != nil {
		q.list.Remove(j.np)
		j.np = nil
	}

	j.Status = job.Concluded
    j.ConcludedAt = time.Now()
	q.value.Store(j.Id, j)
	return nil
}

func (q *Queue) Cancel(id uint64) error {
	jAny, ok := q.value.Load(id)
	if !ok {
		return fmt.Errorf("%w job with id=%d", ErrNotFound, id)
	}
	j := jAny.(jobNP)

    if j.Status == job.Canceled {
        return nil
    }

    if j.Status == job.Concluded {
        return fmt.Errorf("%w job with id=%d was already Concluded", ErrFinalState, id)
    }

	if j.np != nil {
		q.list.Remove(j.np)
		j.np = nil
	}

	j.Status = job.Canceled
	q.value.Store(j.Id, j)
	return nil
}

func (q *Queue) Job(id uint64) (job.Job, error) {
	jAny, ok := q.value.Load(id)
	if !ok {
		return job.Job{}, fmt.Errorf("%w job with id=%d", ErrNotFound, id)
	}
	j := jAny.(jobNP)

	return j.Job, nil
}

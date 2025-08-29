package main

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Repository defines data access behavior (single responsibility for persistence).
type Repository interface {
	Create(ctx context.Context, t *Task) (int64, error)
	Get(ctx context.Context, id int64) (*Task, error)
	Update(ctx context.Context, t *Task) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filter ListFilter) ([]*Task, int, error)
}

var ErrNotFound = errors.New("not found")

// ListFilter for pagination + filtering
type ListFilter struct {
	Status *TaskStatus
	Offset int
	Limit  int
}

type memoryRepo struct {
	mtx     sync.RWMutex
	tasks   map[int64]*Task
	nextID  int64
	nowFunc func() time.Time
}

func NewMemoryRepo() Repository {
	return &memoryRepo{
		tasks:   make(map[int64]*Task),
		nextID:  0,
		nowFunc: time.Now,
	}
}

func (r *memoryRepo) Create(ctx context.Context, t *Task) (int64, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.nextID++
	id := r.nextID
	now := r.nowFunc()
	t.ID = id
	t.CreatedAt = now
	t.UpdatedAt = now
	// clone to avoid outside modification
	c := *t
	r.tasks[id] = &c
	return id, nil
}

func (r *memoryRepo) Get(ctx context.Context, id int64) (*Task, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	t, ok := r.tasks[id]
	if !ok {
		return nil, ErrNotFound
	}
	c := *t
	return &c, nil
}

func (r *memoryRepo) Update(ctx context.Context, t *Task) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	existing, ok := r.tasks[t.ID]
	if !ok {
		return ErrNotFound
	}
	existing.Title = t.Title
	existing.Description = t.Description
	existing.Status = t.Status
	existing.UpdatedAt = r.nowFunc()
	return nil
}

func (r *memoryRepo) Delete(ctx context.Context, id int64) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if _, ok := r.tasks[id]; !ok {
		return ErrNotFound
	}
	delete(r.tasks, id)
	return nil
}

func (r *memoryRepo) List(ctx context.Context, filter ListFilter) ([]*Task, int, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	var all []*Task
	for _, t := range r.tasks {
		if filter.Status != nil && t.Status != *filter.Status {
			continue
		}
		c := *t
		all = append(all, &c)
	}
	total := len(all)

	// bounds check offset/limit
	start := filter.Offset
	if start > total {
		start = total
	}
	end := start + filter.Limit
	if end > total {
		end = total
	}

	return all[start:end], total, nil
}

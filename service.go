package main

import (
	"context"
	"fmt"
)

type TaskService struct {
	repo Repository
}

func NewTaskService(r Repository) *TaskService {
	return &TaskService{repo: r}
}

func (s *TaskService) CreateTask(ctx context.Context, title, desc string) (*Task, error) {
	if title == "" {
		return nil, fmt.Errorf("title required")
	}
	t := &Task{
		Title:       title,
		Description: desc,
		Status:      StatusPending,
	}
	id, err := s.repo.Create(ctx, t)
	if err != nil {
		return nil, err
	}
	return s.repo.Get(ctx, id)
}

func (s *TaskService) GetTask(ctx context.Context, id int64) (*Task, error) {
	return s.repo.Get(ctx, id)
}

func (s *TaskService) UpdateTask(ctx context.Context, id int64, title, desc string, status TaskStatus) (*Task, error) {
	t, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if title != "" {
		t.Title = title
	}
	t.Description = desc
	t.Status = status

	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return s.repo.Get(ctx, id)
}

func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *TaskService) ListTasks(ctx context.Context, page, size int, status *TaskStatus) ([]*Task, int, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	offset := (page - 1) * size
	filter := ListFilter{
		Status: status,
		Offset: offset,
		Limit:  size,
	}
	return s.repo.List(ctx, filter)
}

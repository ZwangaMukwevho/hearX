// pkg/service/task_service.go
package service

import (
	"context"

	"hearx/pkg/model"
	"hearx/pkg/repository"

	"go.uber.org/zap"
)

type TaskService interface {
	AddTask(ctx context.Context, task model.Task) (model.Task, error)
	CompleteTask(ctx context.Context, id int64) (model.Task, error)
	ListTasks(ctx context.Context) ([]model.Task, error)
}

type taskService struct {
	repo   repository.TaskRepository
	logger *zap.Logger
}

func NewTaskService(repo repository.TaskRepository, logger *zap.Logger) TaskService {
	return &taskService{repo: repo, logger: logger}
}

func (s *taskService) AddTask(ctx context.Context, task model.Task) (model.Task, error) {
	s.logger.Info("service: adding task", zap.String("title", task.Title))
	created, err := s.repo.Create(ctx, task)
	if err != nil {
		s.logger.Error("service: AddTask failed", zap.Error(err), zap.Any("task", task))
		return model.Task{}, err
	}
	s.logger.Info("service: task added", zap.Int64("id", created.ID))
	return created, nil
}

func (s *taskService) CompleteTask(ctx context.Context, id int64) (model.Task, error) {
	s.logger.Info("service: completing task", zap.Int64("id", id))
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("service: FindByID failed", zap.Error(err), zap.Int64("id", id))
		return model.Task{}, err
	}
	t.Completed = true
	updated, err := s.repo.Update(ctx, t)
	if err != nil {
		s.logger.Error("service: CompleteTask failed", zap.Error(err), zap.Int64("id", id))
		return model.Task{}, err
	}
	s.logger.Info("service: task completed", zap.Int64("id", updated.ID))
	return updated, nil
}

func (s *taskService) ListTasks(ctx context.Context) ([]model.Task, error) {
	s.logger.Info("service: listing tasks")
	list, err := s.repo.FindAll(ctx)
	if err != nil {
		s.logger.Error("service: ListTasks failed", zap.Error(err))
	}
	return list, err
}

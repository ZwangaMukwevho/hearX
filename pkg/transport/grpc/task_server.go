package grpc

import (
	"context"

	"hearx/pkg/model"
	"hearx/pkg/service"
	pb "hearx/proto/todo"
)

// TaskServer implements the gRPC TodoService.
type TaskServer struct {
	pb.UnimplementedTodoServiceServer
	svc service.TaskService
}

// NewTaskServer constructs a TaskServer with the given business‐logic service.
func NewTaskServer(svc service.TaskService) *TaskServer {
	return &TaskServer{svc: svc}
}

// AddTask creates a new task via the service layer.
func (s *TaskServer) AddTask(ctx context.Context, req *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	// map from proto → internal model
	in := model.Task{
		Title:       req.Task.Title,
		Description: req.Task.Description,
	}

	created, err := s.svc.AddTask(ctx, in)
	if err != nil {
		return nil, err
	}

	// map from internal model → proto
	return &pb.AddTaskResponse{
		Task: &pb.Task{
			Id:          created.ID,
			Title:       created.Title,
			Description: created.Description,
			Completed:   created.Completed,
		},
	}, nil
}

// CompleteTask marks the given task as completed.
func (s *TaskServer) CompleteTask(ctx context.Context, req *pb.CompleteTaskRequest) (*pb.CompleteTaskResponse, error) {
	updated, err := s.svc.CompleteTask(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.CompleteTaskResponse{
		Task: &pb.Task{
			Id:          updated.ID,
			Title:       updated.Title,
			Description: updated.Description,
			Completed:   updated.Completed,
		},
	}, nil
}

// ListTasks retrieves all tasks.
func (s *TaskServer) ListTasks(ctx context.Context, _ *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	list, err := s.svc.ListTasks(ctx)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListTasksResponse{}
	for _, t := range list {
		resp.Tasks = append(resp.Tasks, &pb.Task{
			Id:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Completed:   t.Completed,
		})
	}
	return resp, nil
}

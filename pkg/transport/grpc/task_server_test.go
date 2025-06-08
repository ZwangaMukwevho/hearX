// pkg/transport/grpc/task_server_test.go
package grpc_test

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"hearx/pkg/model"
	mocksvc "hearx/pkg/service/mock_service"
	grpcTransport "hearx/pkg/transport/grpc"
	pb "hearx/proto"
)

var _ = Describe("TaskServer (gRPC)", func() {
	var (
		ctrl    *gomock.Controller
		svcMock *mocksvc.MockTaskService
		server  *grpcTransport.TaskServer
		ctx     context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		svcMock = mocksvc.NewMockTaskService(ctrl)
		server = grpcTransport.NewTaskServer(svcMock)
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("AddTask", func() {
		It("should call service.AddTask and return a mapped response", func() {
			req := &pb.AddTaskRequest{
				Task: &pb.Task{
					Title:       "t1",
					Description: "d1",
				},
			}
			created := model.Task{
				ID:          7,
				Title:       "t1",
				Description: "d1",
				Completed:   false,
			}

			svcMock.
				EXPECT().
				AddTask(ctx, model.Task{Title: "t1", Description: "d1"}).
				Return(created, nil)

			resp, err := server.AddTask(ctx, req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Task.Id).To(Equal(int64(7)))
			Expect(resp.Task.Title).To(Equal("t1"))
			Expect(resp.Task.Completed).To(BeFalse())
		})

		It("should propagate service errors", func() {
			req := &pb.AddTaskRequest{Task: &pb.Task{Title: "oops"}}

			svcMock.
				EXPECT().
				AddTask(ctx, gomock.Any()).
				Return(model.Task{}, errors.New("boom"))

			_, err := server.AddTask(ctx, req)
			Expect(err).To(MatchError("boom"))
		})
	})

	Describe("ListTasks", func() {
		It("should call service.ListTasks and map the result", func() {
			tasks := []model.Task{
				{ID: 1, Title: "A", Completed: false},
				{ID: 2, Title: "B", Completed: true},
			}

			svcMock.
				EXPECT().
				ListTasks(ctx).
				Return(tasks, nil)

			resp, err := server.ListTasks(ctx, &pb.ListTasksRequest{})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(resp.Tasks)).To(Equal(2))
			Expect(resp.Tasks[1].Id).To(Equal(int64(2)))
			Expect(resp.Tasks[1].Completed).To(BeTrue())
		})

		It("should propagate service errors", func() {
			svcMock.
				EXPECT().
				ListTasks(ctx).
				Return(nil, errors.New("fail"))

			_, err := server.ListTasks(ctx, &pb.ListTasksRequest{})
			Expect(err).To(MatchError("fail"))
		})
	})

	Describe("CompleteTask", func() {
		It("should call service.CompleteTask and return a mapped response", func() {
			id := int64(42)
			updated := model.Task{ID: id, Completed: true}

			svcMock.
				EXPECT().
				CompleteTask(ctx, id).
				Return(updated, nil)

			resp, err := server.CompleteTask(ctx, &pb.CompleteTaskRequest{Id: id})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Task.Id).To(Equal(id))
			Expect(resp.Task.Completed).To(BeTrue())
		})

		It("should propagate service errors", func() {
			svcMock.
				EXPECT().
				CompleteTask(ctx, int64(99)).
				Return(model.Task{}, errors.New("nop"))

			_, err := server.CompleteTask(ctx, &pb.CompleteTaskRequest{Id: 99})
			Expect(err).To(MatchError("nop"))
		})
	})
})

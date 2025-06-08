// pkg/service/task_service_test.go
package service_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"

	"hearx/pkg/model"
	mockrepo "hearx/pkg/repository/mock_repository"
	svc "hearx/pkg/service"
)

var _ = Describe("taskService", func() {
	var (
		ctrl     *gomock.Controller
		repoMock *mockrepo.MockTaskRepository
		service  svc.TaskService
		logger   *zap.Logger
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		repoMock = mockrepo.NewMockTaskRepository(ctrl)
		logger = zap.NewNop()
		service = svc.NewTaskService(repoMock, logger)
	})

	AfterEach(func() { ctrl.Finish() })

	Describe("AddTask", func() {
		It("should create and return the new task", func() {
			in := model.Task{Title: "T1", Description: "D1"}
			out := model.Task{ID: 42, Title: "T1", Description: "D1"}

			repoMock.
				EXPECT().
				Create(gomock.Any(), in).
				Return(out, nil)

			result, err := service.AddTask(context.Background(), in)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(out))
		})

		It("should propagate repository errors", func() {
			in := model.Task{Title: "X"}

			repoMock.
				EXPECT().
				Create(gomock.Any(), in).
				Return(model.Task{}, errors.New("boom"))

			_, err := service.AddTask(context.Background(), in)
			Expect(err).To(MatchError("boom"))
		})
	})

	Describe("ListTasks", func() {
		It("should return list from repo", func() {
			list := []model.Task{{ID: 1, Title: "A"}}

			repoMock.
				EXPECT().
				FindAll(gomock.Any()).
				Return(list, nil)

			result, err := service.ListTasks(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(list))
		})

		It("should propagate errors", func() {
			repoMock.
				EXPECT().
				FindAll(gomock.Any()).
				Return(nil, errors.New("fail"))

			_, err := service.ListTasks(context.Background())
			Expect(err).To(MatchError("fail"))
		})
	})

	Describe("CompleteTask", func() {
		It("should fetch, update and return the completed task", func() {
			id := int64(7)
			orig := model.Task{ID: id, Title: "T", Completed: false}
			updated := model.Task{ID: id, Title: "T", Completed: true}

			gomock.InOrder(
				repoMock.EXPECT().FindByID(gomock.Any(), id).Return(orig, nil),
				repoMock.EXPECT().Update(gomock.Any(), updated).Return(updated, nil),
			)

			result, err := service.CompleteTask(context.Background(), id)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(updated))
		})

		It("should propagate FindByID error", func() {
			id := int64(9)

			repoMock.
				EXPECT().
				FindByID(gomock.Any(), id).
				Return(model.Task{}, errors.New("missing"))

			_, err := service.CompleteTask(context.Background(), id)
			Expect(err).To(MatchError("missing"))
		})

		It("should propagate Update error", func() {
			id := int64(11)
			orig := model.Task{ID: id, Title: "X"}

			repoMock.
				EXPECT().
				FindByID(gomock.Any(), id).
				Return(orig, nil)
			repoMock.
				EXPECT().
				Update(gomock.Any(), gomock.Any()).
				Return(model.Task{}, errors.New("nope"))

			_, err := service.CompleteTask(context.Background(), id)
			Expect(err).To(MatchError("nope"))
		})
	})
})

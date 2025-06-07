// pkg/cli/cli.go
package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"hearx/pkg/server"
	pb "hearx/proto"
)

var (
	Host string
	Port string
)

// ServeCmd starts your gRPC server via Fx
func ServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "serve",
		Short:         "Run gRPC server",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// blocks until shutdown
			server.Run()
			return nil
		},
	}
	return cmd
}

// AddCmd calls the AddTask RPC
func AddCmd() *cobra.Command {
	var title, desc string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new task",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := grpc.Dial(
				fmt.Sprintf("%s:%s", Host, Port),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithBlock(),
				grpc.WithTimeout(5*time.Second),
			)
			if err != nil {
				return err
			}
			defer conn.Close()

			client := pb.NewTodoServiceClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			res, err := client.AddTask(ctx, &pb.AddTaskRequest{
				Task: &pb.Task{Title: title, Description: desc},
			})
			if err != nil {
				return err
			}
			fmt.Printf("Created task ID=%d\n", res.Task.Id)
			return nil
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "Task title")
	cmd.MarkFlagRequired("title")
	cmd.Flags().StringVar(&desc, "desc", "", "Task description")
	return cmd
}

// GetCmd calls the ListTasks RPC
func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "List all tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := grpc.Dial(
				fmt.Sprintf("%s:%s", Host, Port),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithBlock(),
				grpc.WithTimeout(5*time.Second),
			)
			if err != nil {
				return err
			}
			defer conn.Close()

			client := pb.NewTodoServiceClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			res, err := client.ListTasks(ctx, &pb.ListTasksRequest{})
			if err != nil {
				return err
			}
			for _, t := range res.Tasks {
				fmt.Printf("[%d] %s (completed=%v)\n", t.Id, t.Title, t.Completed)
			}
			return nil
		},
	}
}

// CompleteCmd calls the CompleteTask RPC
func CompleteCmd() *cobra.Command {
	var id int64
	cmd := &cobra.Command{
		Use:   "complete",
		Short: "Mark a task complete",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := grpc.Dial(
				fmt.Sprintf("%s:%s", Host, Port),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithBlock(),
				grpc.WithTimeout(5*time.Second),
			)
			if err != nil {
				return err
			}
			defer conn.Close()

			client := pb.NewTodoServiceClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if _, err := client.CompleteTask(ctx, &pb.CompleteTaskRequest{Id: id}); err != nil {
				return err
			}
			fmt.Printf("Task %d marked complete\n", id)
			return nil
		},
	}
	cmd.Flags().Int64Var(&id, "id", 0, "Task ID")
	cmd.MarkFlagRequired("id")
	return cmd
}

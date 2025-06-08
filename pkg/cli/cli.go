// pkg/cli/cli.go
package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"hearx/pkg/logger"
	"hearx/pkg/server"
	pb "hearx/proto/todo"
)

var (
	Host string
	Port string
)

func Execute() error {
	root := &cobra.Command{
		Use:   "todo",
		Short: "Todo CLI (serve + client)",
	}

	// global flags for client
	root.PersistentFlags().StringVar(&Host, "host", "localhost", "gRPC host")
	root.PersistentFlags().StringVar(&Port, "port", "50051", "gRPC port")

	root.AddCommand(
		serveCmd(),
		addCmd(),
		getCmd(),
		completeCmd(),
	)
	return root.Execute()
}

func serveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start both gRPC & HTTP-Gateway",
		RunE: func(cmd *cobra.Command, args []string) error {
			// We initialize a logger so we can see startup logs in the CLI
			log, err := logger.NewLogger()
			if err != nil {
				fmt.Fprintf(os.Stderr, "logger init failed: %v\n", err)
				return err
			}
			defer log.Sync()
			log.Info("serve command invoked")

			server.Run()
			return nil
		},
	}
}

func addCmd() *cobra.Command {
	var title, desc string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new task",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := dial()
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
	cmd.Flags().StringVar(&title, "title", "", "Task title (required)")
	cmd.MarkFlagRequired("title")
	cmd.Flags().StringVar(&desc, "desc", "", "Task description")
	return cmd
}

func getCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "List all tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := dial()
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

func completeCmd() *cobra.Command {
	var id int64
	cmd := &cobra.Command{
		Use:   "complete",
		Short: "Mark a task complete",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := dial()
			if err != nil {
				return err
			}
			defer conn.Close()

			client := pb.NewTodoServiceClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err = client.CompleteTask(ctx, &pb.CompleteTaskRequest{Id: id})
			if err != nil {
				return err
			}
			fmt.Printf("Task %d marked complete\n", id)
			return nil
		},
	}
	cmd.Flags().Int64Var(&id, "id", 0, "ID of the task to complete (required)")
	cmd.MarkFlagRequired("id")
	return cmd
}

func dial() (*grpc.ClientConn, error) {
	addr := fmt.Sprintf("%s:%s", Host, Port)
	return grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
}

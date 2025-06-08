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

	"hearx/pkg/server"
	pb "hearx/proto"
)

var (
	// Server flags
	FlagGRPCPort  string
	FlagMySQLHost string
	FlagMySQLPort string
	FlagMySQLUser string
	FlagMySQLPass string
	FlagMySQLDB   string

	// Client flags
	ClientHost string
	ClientPort string
)

func Execute() error {
	root := &cobra.Command{
		Use:   "todo",
		Short: "Todo application CLI (server + client)",
	}

	// server subcommand
	root.AddCommand(serverCmd())

	// client subcommand with its own subcommands
	root.AddCommand(clientCmd())

	return root.Execute()
}

// serverCmd configures and launches your gRPC+Gateway server.
func serverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run the gRPC server (and HTTP-Gateway)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// export flags into env for server.Run to pick up
			os.Setenv("GRPC_PORT", FlagGRPCPort)
			os.Setenv("MYSQL_HOST", FlagMySQLHost)
			os.Setenv("MYSQL_PORT", FlagMySQLPort)
			os.Setenv("MYSQL_USER", FlagMySQLUser)
			os.Setenv("MYSQL_PASSWORD", FlagMySQLPass)
			os.Setenv("MYSQL_DATABASE", FlagMySQLDB)

			// this will block until the process is terminated
			server.Run()
			return nil
		},
	}

	// gRPC listener port
	cmd.Flags().StringVar(&FlagGRPCPort, "grpc-port", "50051", "gRPC listen port")

	// MySQL connection flags
	cmd.Flags().StringVar(&FlagMySQLHost, "mysql-host", "localhost", "MySQL host")
	cmd.Flags().StringVar(&FlagMySQLPort, "mysql-port", "3306", "MySQL port")
	cmd.Flags().StringVar(&FlagMySQLUser, "mysql-user", "user", "MySQL user")
	cmd.Flags().StringVar(&FlagMySQLPass, "mysql-pass", "password", "MySQL password")
	cmd.Flags().StringVar(&FlagMySQLDB, "mysql-db", "project_db", "MySQL database name")

	return cmd
}

// clientCmd groups the client subcommands
func clientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "Run the gRPC client (add|get|complete)",
	}

	// client flags apply to all subcommands
	cmd.PersistentFlags().StringVar(&ClientHost, "host", "localhost", "gRPC server host")
	cmd.PersistentFlags().StringVar(&ClientPort, "port", "50051", "gRPC server port")

	// add the actions
	cmd.AddCommand(addCmd())
	cmd.AddCommand(getCmd())
	cmd.AddCommand(completeCmd())
	return cmd
}

// addCmd calls the AddTask RPC
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

// getCmd calls the ListTasks RPC
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

// completeCmd calls the CompleteTask RPC
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

			if _, err := client.CompleteTask(ctx, &pb.CompleteTaskRequest{Id: id}); err != nil {
				return err
			}
			fmt.Printf("Task %d marked complete\n", id)
			return nil
		},
	}
	cmd.Flags().Int64Var(&id, "id", 0, "Task ID (required)")
	cmd.MarkFlagRequired("id")
	return cmd
}

// dial creates a gRPC connection to Host:Port.
func dial() (*grpc.ClientConn, error) {
	addr := fmt.Sprintf("%s:%s", ClientHost, ClientPort)
	return grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
}

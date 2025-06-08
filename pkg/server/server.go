// pkg/server/server.go
package server

import (
	"context"
	"fmt"
	"net"
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"hearx/pkg/auth"
	"hearx/pkg/logger"
	"hearx/pkg/repository"
	"hearx/pkg/service"
	"hearx/pkg/storage"
	grpcTransport "hearx/pkg/transport/grpc"
	pb "hearx/proto"
)

func Run() {
	app := fx.New(
		fx.Provide(
			logger.NewLogger,
			provideMySQLDSN,
			storage.NewMySQLConn,
			repository.NewTaskRepository,
			service.NewTaskService,
			grpcTransport.NewTaskServer,
			newGRPCServer,
			newListener,
		),
		fx.Invoke(register, start),
	)
	app.Run()
}

func provideMySQLDSN() string {
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	db := os.Getenv("MYSQL_DATABASE")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, pass, host, port, db,
	)
}

func newGRPCServer() *grpc.Server {
	return grpc.NewServer(
		grpc.UnaryInterceptor(auth.UnaryServerInterceptor()),
	)
}
func newListener() (net.Listener, error) {
	p := os.Getenv("GRPC_PORT")
	if p == "" {
		p = "50051"
	}
	return net.Listen("tcp", ":"+p)
}

func register(server *grpc.Server, ts *grpcTransport.TaskServer) {
	pb.RegisterTodoServiceServer(server, ts)
}

func start(lc fx.Lifecycle, server *grpc.Server, lis net.Listener, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("gRPC starting", zap.String("addr", lis.Addr().String()))
			go server.Serve(lis)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("gRPC stopping")
			server.GracefulStop()
			return nil
		},
	})
}

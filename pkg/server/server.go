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

	"hearx/pkg/logger"
	"hearx/pkg/repository"
	"hearx/pkg/service"
	"hearx/pkg/storage"
	grpcTransport "hearx/pkg/transport/grpc"
	pb "hearx/proto/todo"
)

// Run boots both gRPC and HTTP-Gateway in one process.
func Run() {
	app := fx.New(
		fx.Provide(
			// common
			logger.NewLogger,
			provideMySQLDSN,
			storage.NewMySQLConn,
			repository.NewTaskRepository,
			service.NewTaskService,

			// --- gRPC server providers ---
			grpcTransport.NewTaskServer,
			newGRPCServer,
			fx.Annotate(
				newGRPCListener,
				fx.ResultTags(`name:"grpcListener"`),
			),

			// --- HTTP-Gateway providers ---
			NewGatewayMux,
			fx.Annotate(
				NewGatewayListener,
				fx.ResultTags(`name:"httpListener"`),
			),
		),
		fx.Invoke(
			// start gRPC
			registerGRPC,
			fx.Annotate(
				startGRPC,
				fx.ParamTags(``, ``, `name:"grpcListener"`, ``),
			),

			// start HTTP-Gateway
			fx.Annotate(
				startGateway,
				fx.ParamTags(``, ``, `name:"httpListener"`, ``),
			),
		),
	)
	app.Run()
}

// build DSN from env
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

// ─── gRPC setup ────────────────────────────────────────────────

func newGRPCServer() *grpc.Server {
	return grpc.NewServer()
}

func newGRPCListener() (net.Listener, error) {
	p := os.Getenv("GRPC_PORT")
	if p == "" {
		p = "50051"
	}
	return net.Listen("tcp", ":"+p)
}

func registerGRPC(s *grpc.Server, ts *grpcTransport.TaskServer) {
	pb.RegisterTodoServiceServer(s, ts)
}

func startGRPC(lc fx.Lifecycle, s *grpc.Server, lis net.Listener, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("gRPC starting", zap.String("addr", lis.Addr().String()))
			go s.Serve(lis)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("gRPC stopping")
			s.GracefulStop()
			return nil
		},
	})
}

// ─── HTTP-Gateway setup ────────────────────────────────────────

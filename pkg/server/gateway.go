// pkg/server/gateway.go
package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "hearx/proto/todo"
)

// NewGatewayMux constructs the HTTP→gRPC mux.
func NewGatewayMux() *runtime.ServeMux {
	return runtime.NewServeMux()
}

// NewGatewayListener listens on HTTP_PORT (default 8000).
func NewGatewayListener() (net.Listener, error) {
	p := os.Getenv("HTTP_PORT")
	if p == "" {
		p = "8000"
	}
	return net.Listen("tcp", ":"+p)
}

// startGateway registers the handlers and starts HTTP.
func startGateway(lc fx.Lifecycle, mux *runtime.ServeMux, lis net.Listener, log *zap.Logger) {
	srv := &http.Server{Addr: lis.Addr().String(), Handler: mux}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// wire generated handlers to the mux
			endpoint := fmt.Sprintf("localhost:%s", os.Getenv("GRPC_PORT"))
			opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
			if err := pb.RegisterTodoServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
				return fmt.Errorf("gateway registration failed: %w", err)
			}
			log.Info("HTTP-Gateway starting", zap.String("addr", lis.Addr().String()))
			go srv.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("HTTP-Gateway stopping")
			return srv.Shutdown(ctx)
		},
	})
}

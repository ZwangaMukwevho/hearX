// pkg/auth/auth.go
package auth

import (
	"context"
	"fmt"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor checks for a valid Bearer token in metadata.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	expected := os.Getenv("AUTH_TOKEN")
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "no metadata")
		}
		vals := md["authorization"]
		if len(vals) == 0 {
			return nil, status.Error(codes.Unauthenticated, "no auth token")
		}
		token := strings.TrimPrefix(vals[0], "Bearer ")
		if token != expected {
			return nil, status.Error(codes.Unauthenticated, "invalid auth token")
		}
		// ok!
		return handler(ctx, req)
	}
}

// StaticTokenCreds implements PerRPCCredentials by always sending the same token.
type StaticTokenCreds struct {
	Token string
}

// GetRequestMetadata injects the authorization header.
func (c StaticTokenCreds) GetRequestMetadata(ctx context.Context, _ ...string) (map[string]string, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("no auth token provided")
	}
	return map[string]string{
		"authorization": "Bearer " + c.Token,
	}, nil
}

// RequireTransportSecurity controls whether this credential requires TLS.
// Return false if you’re on plaintext; for production you’d typically return true.
func (c StaticTokenCreds) RequireTransportSecurity() bool {
	return false
}

// AsPerRPCCreds wraps StaticTokenCreds as oauth.TokenSource so you can also do:
// grpc.WithPerRPCCredentials(oauth.TokenSource{TokenSource: myStaticSource})
func TokenSource(token string) credentials.PerRPCCredentials {
	return StaticTokenCreds{Token: token}
}

package gateway

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type HandlerFromEndpoint func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

// NewGateway creates a new v0.server v0.server for handling HTTP requests.
//
// It takes a context.Context, an endpoint string, and one or more HandlerFromEndpoint functions as parameters.
//
// It returns an http.Handler and an error. The http.Handler is the v0.server v0.server that will handle incoming HTTP requests.
// The error is returned if there are no handlers provided.
func NewGateway(ctx context.Context, endpoint string, handlers ...HandlerFromEndpoint) (*runtime.ServeMux, error) {
	if len(handlers) < 1 {
		return nil, errors.New("NO_HANDLERS")
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	mux := runtime.NewServeMux()
	for _, handler := range handlers {
		if err := handler(ctx, mux, endpoint, opts); err != nil {
			return nil, err
		}
	}

	return mux, nil
}

// GRPCHandlerFunc generates an http.Handler that can be used to handle gRPC requests
// and HTTP requests that are not gRPC requests.
//
// It takes in a *grpc.Server and an http.Handler as parameters.
// The *grpc.Server is used to handle gRPC requests.
// The http.Handler is used to handle non-gRPC HTTP requests.
//
// The function returns an http.Handler that can be used to handle both gRPC and non-gRPC requests.
func GRPCHandlerFunc(grpcServer *grpc.Server, muxgw http.Handler) http.Handler {
	if muxgw == nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			grpcServer.ServeHTTP(w, r)
		})
	}
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			muxgw.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

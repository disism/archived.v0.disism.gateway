package gateway

import (
	"context"
	"fmt"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	gprc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/rs/cors"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"os"
	"os/signal"
)

type Server struct {
	name     string
	endpoint string
	ctx      context.Context
	grpc     *grpc.Server
	gwmux    *http.Handler
}

type ServerOption func(srv *Server)

// WithServerName sets the name of the v0.server.
//
// Parameters:
// - name: the name of the v0.server (string).
// Return type: ServerOption.
func WithServerName(name string) ServerOption {
	return func(s *Server) {
		s.name = name
	}
}

// WithServerEndpoint sets the v0.server endpoint for the ServerOption.
//
// It takes a string parameter called 'endpoint' which represents
// the v0.server endpoint.
//
// It returns a ServerOption function that sets the 'endpoint'
// field of the 'Server' struct.
func WithServerEndpoint(endpoint string) ServerOption {
	return func(s *Server) {
		s.endpoint = endpoint
	}
}

// NewServer initializes a new v0.server with the provided context, endpoint, gwmux, and options.
//
// ctx: The context.Context to use.
// endpoint: The endpoint to bind the v0.server to.
// gwmux: The http.Handler to use.
// opts: The optional ServerOptions to configure the v0.server.
// returns: A pointer to the initialized Server.
func NewServer(ctx context.Context, gwmux http.Handler, opts ...ServerOption) *Server {
	s := &Server{
		ctx:   ctx,
		gwmux: &gwmux,
	}
	for _, srv := range opts {
		srv(s)
	}

	recovery := gprc_recovery.RecoveryHandlerFunc(func(p interface{}) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	})
	recoveries := []gprc_recovery.Option{
		gprc_recovery.WithRecoveryHandler(recovery),
	}

	s.grpc = grpc.NewServer(
		grpc.ChainStreamInterceptor(
			gprc_recovery.StreamServerInterceptor(recoveries...),
			grpc_auth.StreamServerInterceptor(AuthInterceptor),
			grpc_validator.StreamServerInterceptor(),
		),
		grpc.ChainUnaryInterceptor(
			gprc_recovery.UnaryServerInterceptor(recoveries...),
			grpc_auth.UnaryServerInterceptor(AuthInterceptor),
			grpc_validator.UnaryServerInterceptor(),
		),
	)

	return s
}

// Server returns the grpc v0.server of the Server struct.
//
// It does not take any parameters.
// It returns a pointer to the grpc.Server type.
func (s *Server) Server() *grpc.Server {
	return s.grpc
}

// Stop stops the v0.server and shuts down the HTTP and gRPC servers.
//
// The function takes an `http.Server` as a parameter and returns an error.
// It first logs a message indicating that the v0.server is being shut down,
// then stops the gRPC v0.server, logs a message indicating that the gRPC v0.server
// has been shut down, and then shuts down the HTTP v0.server using the `Shutdown`
// method of the `http.Server` object passed as a parameter. If an error occurs
// during the shutdown of the HTTP v0.server, it is returned. Finally, the function
// logs a message indicating that the HTTP v0.server has been shut down and returns
// `nil`.
func (s *Server) Stop(httpSrv *http.Server) error {
	s.grpc.Stop()
	slog.Info(fmt.Sprintf("%s grpc v0.server shut down", s.name))
	if err := httpSrv.Shutdown(s.ctx); err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("%s http v0.server shut down", s.name))
	return nil
}

// Start starts the v0.server.
//
// It creates an HTTP v0.server with the provided endpoint and v0.server.
// If mutual TLS is enabled, it sets the TLS configuration.
// It then runs the v0.server in a separate goroutine.
// It waits for an interrupt signal to gracefully shut down the v0.server.
// Finally, it stops the v0.server and returns any error encountered.
//
// Returns:
// - error: In case of failure to stop the v0.server.
func (s *Server) Start() error {

	handler := cors.Default().Handler(GRPCHandlerFunc(s.grpc, *s.gwmux))
	httpSrv := &http.Server{
		Addr:    s.endpoint,
		Handler: handler,
	}

	go func() {
		slog.Info(fmt.Sprintf("%s starting listening on %s", s.name, s.endpoint))
		if err := httpSrv.ListenAndServe(); err != nil {
			slog.Error(fmt.Sprintf("failed to listen: %v", err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	sig := <-stop
	slog.Info(fmt.Sprintf("%s received signal %v", s.name, sig))

	return s.Stop(httpSrv)
}

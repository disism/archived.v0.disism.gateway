package gateway

import (
	"github.com/disism/v0.disism.account/sdk"
	"fmt"
	"gateway/jwt"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ignore = []string{
	"/proto.account.v1alpha1.AuthenticationService/Authentication",
	"/proto.account.v1alpha1.UserService/CreateUser",
	"/proto.account.v1alpha1.UserService/GetUser",
	"/proto.account.v1alpha1.DeviceService/GetDevice",
	"/proto.account.v1alpha1.PublicService/GetWebfinger",
	"/proto.account.v1alpha1.PublicService/GetVersion",
}

const (
	JWTMetadata = "jwt_meta"
)

func AuthInterceptor(ctx context.Context) (context.Context, error) {
	method, _ := grpc.Method(ctx)
	for _, i := range ignore {
		if method == i {
			return ctx, nil
		}
	}
	fmt.Println(method)

	b, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, err
	}

	parse, err := jwt.JWTParse(b)
	if err != nil {
		return nil, err
	}

	d, err := sdk.GetDevice(ctx, parse.ID)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, codes.Unauthenticated.String())
	}
	return context.WithValue(ctx, JWTMetadata, &Userdata{
		ID:       d.GetUser().GetId(),
		Username: d.GetUser().GetUsername(),
		DeviceID: d.GetDevice().GetId(),
		IP:       d.GetDevice().GetIp(),
	}), nil
}

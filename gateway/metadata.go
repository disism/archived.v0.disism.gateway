package gateway

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

const (
	gatewayUserAgent = "grpcgateway-user-agent"
	forwarded        = "x-forwarded-for"
)

type Userdata struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	DeviceID uint64 `json:"device"`
	IP       string `json:"ip"`
}

type Matadate struct {
	ctx context.Context
	Userdata
}

func NewMatadata(ctx context.Context) *Matadate {
	md := ctx.Value(JWTMetadata).(*Userdata)
	return &Matadate{
		ctx: ctx,
		Userdata: Userdata{
			ID:       md.ID,
			Username: md.Username,
			DeviceID: md.DeviceID,
			IP:       md.IP,
		},
	}
}

func (md *Matadate) GetUserID() (ID uint64) {
	return md.ID
}

func (md *Matadate) GetUsername() (username string) {
	return md.Username
}

func (md *Matadate) GetDeviceID() (deviceID uint64) {
	return md.DeviceID
}

func (md *Matadate) GetUserIP() (ip string) {
	return md.IP
}

func GetUAFromContext(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	ua := md.Get(gatewayUserAgent)
	return ua[0]
}

func GetIPFromContext(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	ip := md.Get(forwarded)
	return ip[0]
}

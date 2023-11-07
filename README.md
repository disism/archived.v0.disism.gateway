# Archived!

RUN SERVER EXAMPLE...
```
type server struct {
	auth.AuthServer
	auth.AccountServer
	auth.OAuth2Server
}

func Run(endpoint string) error {
	ctx := context.Background()

	gateway, err := gateway.NewGateway(
	    ctx,
	    endpoint,
		auth.RegisterAccountHandlerFromEndpoint,
	)
	if err != nil {
		return err
	}
	
	srv := gateway.NewServer(ctx, endpoint, gateway)
	auth.RegisterAccountServer(srv.GetGRPCServer(), &server{})
	auth.RegisterAuthServer(srv.GetGRPCServer(), &server{})

	if err := srv.Start(); err != nil {
		return err
	}
	return nil
}


```



CLIENT
```
conn, err := grpc.Dial(URL, )
if err != nil {
    log.Fatalf("did not connect: %v", err)
}
```
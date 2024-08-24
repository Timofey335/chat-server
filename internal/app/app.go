package app

import (
	"context"
	"log"
	"net"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/Timofey335/chat-server/internal/closer"
	"github.com/Timofey335/chat-server/internal/config"
	desc "github.com/Timofey335/chat-server/pkg/chat_server_v1"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

func NewApp(ctx context.Context, cfg string) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx, cfg); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGRPCServer()
}

func (a *App) initDeps(ctx context.Context, cfg string) error {

	inits := []func(context.Context, string) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		if err := f(ctx, cfg); err != nil {
			return err
		}

	}

	return nil
}

func (a *App) initConfig(_ context.Context, cfg string) error {
	if err := config.Load(cfg); err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context, _ string) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context, _ string) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	desc.RegisterChatServerV1Server(a.grpcServer, a.serviceProvider.ServImplementation(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf(color.BlueString("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address()))

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	if err = a.grpcServer.Serve(list); err != nil {
		return err
	}

	return nil
}

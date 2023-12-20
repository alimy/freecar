package main

import (
	"context"

	"github.com/alimy/freecar/app/user/config"
	"github.com/alimy/freecar/app/user/internal"
	"github.com/alimy/freecar/app/user/rpc"
	"github.com/alimy/freecar/app/user/servants"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
)

func main() {
	// initialization
	internal.Initial()
	rpc.Initial()

	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	// Create new server.
	srv := servants.NewUserService()
	if err := srv.Run(); err != nil {
		klog.Fatal(err)
	}
}

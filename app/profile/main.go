package main

import (
	"context"

	"github.com/alimy/freecar/app/profile/config"
	"github.com/alimy/freecar/app/profile/internal"
	"github.com/alimy/freecar/app/profile/servants"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
)

func main() {
	// initialization
	internal.Initial()

	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	// Create new server.
	srv := servants.NewProfileService()
	if err := srv.Run(); err != nil {
		klog.Fatal(err)
	}
}

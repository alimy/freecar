package main

import (
	"context"

	"github.com/alimy/freecar/app/blob/conf"
	"github.com/alimy/freecar/app/blob/internal"
	"github.com/alimy/freecar/app/blob/servants"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
)

func main() {
	// initialization
	conf.Initial()
	internal.Initial()

	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(conf.GlobalServerConfig.Name),
		provider.WithExportEndpoint(conf.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	// Create new server.
	srv := servants.NewBlobService()
	if err := srv.Run(); err != nil {
		klog.Fatal(err)
	}
}

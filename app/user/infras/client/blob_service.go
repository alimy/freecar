package client

import (
	"fmt"

	"github.com/alimy/freecar/app/user/conf"
	"github.com/alimy/freecar/idle/auto/rpc/blob/blobservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
)

// initBlob to init blob service
func initBlob() blobservice.Client {
	// init resolver
	r, err := consul.NewConsulResolver(fmt.Sprintf("%s:%d",
		conf.GlobalConsulConfig.Host,
		conf.GlobalConsulConfig.Port))
	if err != nil {
		klog.Fatalf("new nacos client failed: %s", err.Error())
	}
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(conf.GlobalServerConfig.BlobSrvInfo.Name),
		provider.WithExportEndpoint(conf.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)

	// create a new client
	c, err := blobservice.NewClient(
		conf.GlobalServerConfig.BlobSrvInfo.Name,
		client.WithResolver(r),                                     // service discovery
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()), // load balance
		client.WithMuxConnection(1),                                // multiplexing
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: conf.GlobalServerConfig.BlobSrvInfo.Name}),
	)
	if err != nil {
		klog.Fatalf("ERROR: cannot init client: %v\n", err)
	}
	return c
}

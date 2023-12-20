package rpc

import (
	"fmt"

	"github.com/alimy/freecar/app/trip/config"
	"github.com/alimy/freecar/idle/auto/rpc/profile/profileservice"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
)

// initProfile to init profile service
func initProfile() profileservice.Client {
	// init resolver
	r, err := consul.NewConsulResolver(fmt.Sprintf("%s:%d",
		config.GlobalConsulConfig.Host,
		config.GlobalConsulConfig.Port))
	if err != nil {
		hlog.Fatalf("new consul client failed: %s", err.Error())
	}
	// init OpenTelemetry
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.ProfileSrvInfo.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	// create a new client
	c, err := profileservice.NewClient(
		config.GlobalServerConfig.ProfileSrvInfo.Name,
		client.WithResolver(r),                                     // service discovery
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()), // load balance
		client.WithMuxConnection(1),                                // multiplexing
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.GlobalServerConfig.ProfileSrvInfo.Name}),
	)
	if err != nil {
		klog.Fatalf("ERROR: cannot init client: %v\n", err)
	}
	return c
}

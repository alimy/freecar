package rpc

import (
	"fmt"

	"github.com/alimy/freecar/app/car/config"
	"github.com/alimy/freecar/idle/auto/rpc/trip/tripservice"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
)

// initTrip to init trip service
func initTrip() {
	// init resolver
	r, err := consul.NewConsulResolver(fmt.Sprintf("%s:%d",
		config.GlobalConsulConfig.Host,
		config.GlobalConsulConfig.Port))
	if err != nil {
		hlog.Fatalf("new consul client failed: %s", err.Error())
	}
	// init OpenTelemetry
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.TripSrvInfo.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)

	// create a new client
	TripSvc, err = tripservice.NewClient(
		config.GlobalServerConfig.TripSrvInfo.Name,
		client.WithResolver(r),                                     // service discovery
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()), // load balance
		client.WithMuxConnection(1),                                // multiplexing
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.GlobalServerConfig.TripSrvInfo.Name}),
	)
	if err != nil {
		klog.Fatalf("ERROR: cannot init client: %v\n", err)
	}
}

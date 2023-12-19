package main

import (
	"context"
	"net"
	"strconv"

	"github.com/alimy/freecar/app/api/cmd/trip/config"
	"github.com/alimy/freecar/app/api/cmd/trip/initialize"
	"github.com/alimy/freecar/app/api/cmd/trip/pkg/car"
	"github.com/alimy/freecar/app/api/cmd/trip/pkg/mongo"
	"github.com/alimy/freecar/app/api/cmd/trip/pkg/pay"
	"github.com/alimy/freecar/app/api/cmd/trip/pkg/poi"
	"github.com/alimy/freecar/app/api/cmd/trip/pkg/profile"
	"github.com/alimy/freecar/idle/auto/rpc/trip/tripservice"
	"github.com/alimy/freecar/library/cor/consts"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
)

func main() {
	// initialization
	initialize.InitLogger()
	initialize.InitConfig()
	IP, Port := initialize.InitFlag()
	r, info := initialize.InitRegistry(Port)
	db := initialize.InitDB()
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())
	initialize.InitCar()
	initialize.InitProfile()
	initialize.InitUser()

	impl := new(TripServiceImpl)
	impl.PayManager = pay.NewManager(config.UserClient)

	impl.CarManager = &car.Manager{
		CarService: config.CarClient,
	}
	impl.ProfileManager = &profile.Manager{
		ProfileService: config.ProfileClient,
	}
	impl.POIManager = &poi.Manager{}

	impl.MongoManager = mongo.NewManager(db)
	// Create new server.
	srv := tripservice.NewServer(impl,
		server.WithServiceAddr(utils.NewNetAddr(consts.TCP, net.JoinHostPort(IP, strconv.Itoa(Port)))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.GlobalServerConfig.Name}),
	)

	err := srv.Run()
	if err != nil {
		klog.Fatal(err)
	}
}

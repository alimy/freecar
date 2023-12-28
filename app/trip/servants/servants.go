package servants

import (
	"net"
	"strconv"

	"github.com/alimy/freecar/app/trip/conf"
	"github.com/alimy/freecar/app/trip/infras/car"
	"github.com/alimy/freecar/app/trip/infras/client"
	"github.com/alimy/freecar/app/trip/infras/mongo"
	"github.com/alimy/freecar/app/trip/infras/pay"
	"github.com/alimy/freecar/app/trip/infras/profile"
	"github.com/alimy/freecar/app/trip/internal"
	"github.com/alimy/freecar/idle/auto/rpc/trip/tripservice"
	"github.com/alimy/freecar/library/core/consts"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
)

func NewTripService() server.Server {
	ip, port := internal.InitFlag()
	r, info := internal.InitRegistry(port)
	db := internal.NewMongoDB()
	tss := &tripSrv{
		ProfileManager: &profile.Manager{
			ProfileService: client.GetProfileService(),
		},
		CarManager: &car.Manager{
			CarService: client.GetCarService(),
		},
		POIManager:   &poiManager{},
		MongoManager: mongo.NewManager(db),
		PayManager:   pay.NewManager(client.GetUserService()),
	}
	return tripservice.NewServer(tss,
		server.WithServiceAddr(utils.NewNetAddr(consts.TCP, net.JoinHostPort(ip, strconv.Itoa(port)))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: conf.GlobalServerConfig.Name}),
	)
}

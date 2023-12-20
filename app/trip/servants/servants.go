package servants

import (
	"net"
	"strconv"

	"github.com/alimy/freecar/app/trip/config"
	"github.com/alimy/freecar/app/trip/internal"
	"github.com/alimy/freecar/app/trip/pkg/car"
	"github.com/alimy/freecar/app/trip/pkg/mongo"
	"github.com/alimy/freecar/app/trip/pkg/pay"
	"github.com/alimy/freecar/app/trip/pkg/poi"
	"github.com/alimy/freecar/app/trip/pkg/profile"
	"github.com/alimy/freecar/app/trip/rpc"
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
	db := internal.InitDB()
	tss := &tripSrv{
		ProfileManager: &profile.Manager{
			ProfileService: rpc.GetProfileService(),
		},
		CarManager: &car.Manager{
			CarService: rpc.GetCarService(),
		},
		POIManager:   &poi.Manager{},
		MongoManager: mongo.NewManager(db),
		PayManager:   pay.NewManager(rpc.GetUserService()),
	}
	return tripservice.NewServer(tss,
		server.WithServiceAddr(utils.NewNetAddr(consts.TCP, net.JoinHostPort(ip, strconv.Itoa(port)))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.GlobalServerConfig.Name}),
	)
}

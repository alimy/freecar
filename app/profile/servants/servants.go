package servants

import (
	"net"
	"strconv"

	"github.com/alimy/freecar/app/profile/conf"
	"github.com/alimy/freecar/app/profile/infras/client"
	"github.com/alimy/freecar/app/profile/infras/mongo"
	"github.com/alimy/freecar/app/profile/infras/ocr"
	"github.com/alimy/freecar/app/profile/infras/redis"
	"github.com/alimy/freecar/app/profile/internal"
	"github.com/alimy/freecar/idle/auto/rpc/profile/profileservice"
	"github.com/alimy/freecar/library/core/consts"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
)

func NewProfileService() server.Server {
	IP, Port := internal.InitFlag()
	r, info := internal.InitRegistry(Port)
	mongoDb := internal.InitDB()
	redisClient := internal.InitRedis()

	pss := &profileSrv{
		MongoManager:   mongo.NewManager(mongoDb),
		RedisManager:   redis.NewManager(redisClient),
		BlobManager:    client.GetBlobService(),
		LicenseManager: &ocr.LicenseManager{},
	}
	return profileservice.NewServer(pss,
		server.WithServiceAddr(utils.NewNetAddr(consts.TCP, net.JoinHostPort(IP, strconv.Itoa(Port)))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: conf.GlobalServerConfig.Name}),
	)
}

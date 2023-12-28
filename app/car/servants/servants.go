package servants

import (
	"net"
	"strconv"

	"github.com/alimy/freecar/app/car/conf"
	mongoPkg "github.com/alimy/freecar/app/car/infras/mongo"
	"github.com/alimy/freecar/app/car/infras/mq/amqpclt"
	redisPkg "github.com/alimy/freecar/app/car/infras/redis"
	"github.com/alimy/freecar/app/car/internal"
	"github.com/alimy/freecar/idle/auto/rpc/car/carservice"
	"github.com/alimy/freecar/library/core/consts"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
)

// NewCarService create CarService server.
func NewCarService() server.Server {
	ip, port := internal.InitFlag()
	r, info := internal.InitRegistry(port)
	db := internal.InitDB()
	redisClient := internal.InitRedis()
	amqpC := internal.InitMq()
	mqInfo := conf.GlobalServerConfig.RabbitMqInfo
	publisher, err := amqpclt.NewPublisher(amqpC, mqInfo.Exchange)
	if err != nil {
		klog.Fatal("cannot create publisher")
	}
	css := &carSrv{
		Publisher:    publisher,
		MongoManager: mongoPkg.NewManager(db),
		RedisManager: redisPkg.NewManager(redisClient),
	}
	return carservice.NewServer(css,
		server.WithServiceAddr(utils.NewNetAddr(consts.TCP, net.JoinHostPort(ip, strconv.Itoa(port)))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: conf.GlobalServerConfig.Name}),
	)
}

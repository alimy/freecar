package servants

import (
	"net"
	"strconv"

	"github.com/alimy/freecar/app/blob/conf"
	"github.com/alimy/freecar/app/blob/infras/minio"
	"github.com/alimy/freecar/app/blob/infras/mysql"
	"github.com/alimy/freecar/app/blob/infras/redis"
	"github.com/alimy/freecar/app/blob/internal"
	"github.com/alimy/freecar/idle/auto/rpc/blob/blobservice"
	"github.com/alimy/freecar/library/core/consts"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
)

func NewBlobService() server.Server {
	ip, port := internal.InitFlag()
	r, info := internal.InitRegistry(port)
	db := internal.NewDB()
	minioClient := internal.InitMinio()
	redisClient := internal.InitRedis()
	bss := &blobSrv{
		redisManager: redis.NewManager(redisClient),
		minioManager: minio.NewManager(minioClient, conf.GlobalServerConfig.MinioInfo.Bucket),
		mysqlManager: mysql.NewManager(db),
	}
	return blobservice.NewServer(bss,
		server.WithServiceAddr(utils.NewNetAddr(consts.TCP, net.JoinHostPort(ip, strconv.Itoa(port)))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: conf.GlobalServerConfig.Name}),
	)
}

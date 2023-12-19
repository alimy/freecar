package main

import (
	"context"
	"net"
	"strconv"

	"github.com/alimy/freecar/app/api/cmd/blob/config"
	"github.com/alimy/freecar/app/api/cmd/blob/initialize"
	"github.com/alimy/freecar/app/api/cmd/blob/pkg/minio"
	"github.com/alimy/freecar/app/api/cmd/blob/pkg/mysql"
	"github.com/alimy/freecar/app/api/cmd/blob/pkg/redis"
	"github.com/alimy/freecar/idle/auto/rpc/blob/blobservice"
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
	minioClient := initialize.InitMinio()
	redisClient := initialize.InitRedis()
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	// Create new server.
	srv := blobservice.NewServer(&BlobServiceImpl{
		redisManager: redis.NewManager(redisClient),
		minioManager: minio.NewManager(minioClient, config.GlobalServerConfig.MinioInfo.Bucket),
		mysqlManager: mysql.NewManager(db),
	},
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

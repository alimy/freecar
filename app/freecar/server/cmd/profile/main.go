package main

import (
	"context"
	"net"
	"strconv"

	"github.com/alimy/freecar/app/api/cmd/profile/config"
	"github.com/alimy/freecar/app/api/cmd/profile/initialize"
	"github.com/alimy/freecar/app/api/cmd/profile/pkg/mongo"
	"github.com/alimy/freecar/app/api/cmd/profile/pkg/ocr"
	"github.com/alimy/freecar/app/api/cmd/profile/pkg/redis"
	"github.com/alimy/freecar/idle/auto/rpc/profile/profileservice"
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
	mongoDb := initialize.InitDB()
	redisClient := initialize.InitRedis()
	blobClient := initialize.InitBlob()
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	// Create new server.
	srv := profileservice.NewServer(&ProfileServiceImpl{
		MongoManager:   mongo.NewManager(mongoDb),
		RedisManager:   redis.NewManager(redisClient),
		BlobManager:    blobClient,
		LicenseManager: &ocr.LicenseManager{},
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

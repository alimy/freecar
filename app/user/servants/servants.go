package servants

import (
	"net"
	"strconv"

	"github.com/alimy/freecar/app/user/config"
	"github.com/alimy/freecar/app/user/internal"
	"github.com/alimy/freecar/app/user/pkg/md5"
	"github.com/alimy/freecar/app/user/pkg/mysql"
	"github.com/alimy/freecar/app/user/pkg/paseto"
	"github.com/alimy/freecar/app/user/pkg/wechat"
	"github.com/alimy/freecar/app/user/rpc"
	"github.com/alimy/freecar/idle/auto/rpc/user/userservice"
	"github.com/alimy/freecar/library/core/consts"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
)

// NewUserService create UserService server.
func NewUserService() server.Server {
	ip, port := internal.InitFlag()
	r, info := internal.InitRegistry(port)
	db := internal.InitDB()
	tg, err := paseto.NewTokenGenerator(
		config.GlobalServerConfig.PasetoInfo.SecretKey,
		[]byte(config.GlobalServerConfig.PasetoInfo.Implicit))
	if err != nil {
		klog.Fatal(err)
	}
	uss := &userSrv{
		OpenIDResolver: &wechat.AuthServiceImpl{
			AppID:     config.GlobalServerConfig.WXInfo.AppId,
			AppSecret: config.GlobalServerConfig.WXInfo.AppSecret,
		},
		EncryptManager:    &md5.EncryptManager{Salt: config.GlobalServerConfig.MysqlInfo.Salt},
		UserMysqlManager:  mysql.NewUserManager(db, config.GlobalServerConfig.MysqlInfo.Salt),
		AdminMysqlManager: mysql.NewAdminManager(db, config.GlobalServerConfig.MysqlInfo.Salt),
		BlobManager:       rpc.BlobSvc,
		TokenGenerator:    tg,
	}
	return userservice.NewServer(uss,
		server.WithServiceAddr(utils.NewNetAddr(consts.TCP, net.JoinHostPort(ip, strconv.Itoa(port)))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.GlobalServerConfig.Name}),
	)
}

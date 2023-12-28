package servants

import (
	"net"
	"strconv"

	"github.com/alimy/freecar/app/user/conf"
	"github.com/alimy/freecar/app/user/infras/client"
	"github.com/alimy/freecar/app/user/infras/mysql"
	"github.com/alimy/freecar/app/user/internal"
	"github.com/alimy/freecar/app/user/internal/md5"
	"github.com/alimy/freecar/app/user/internal/paseto"
	"github.com/alimy/freecar/app/user/internal/wechat"
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
	db := internal.NewDB()
	tg, err := paseto.NewTokenGenerator(
		conf.GlobalServerConfig.PasetoInfo.SecretKey,
		[]byte(conf.GlobalServerConfig.PasetoInfo.Implicit))
	if err != nil {
		klog.Fatal(err)
	}
	uss := &userSrv{
		OpenIDResolver: &wechat.AuthServiceImpl{
			AppID:     conf.GlobalServerConfig.WXInfo.AppId,
			AppSecret: conf.GlobalServerConfig.WXInfo.AppSecret,
		},
		EncryptManager:    &md5.EncryptManager{Salt: conf.GlobalServerConfig.MysqlInfo.Salt},
		UserMysqlManager:  mysql.NewUserManager(db, conf.GlobalServerConfig.MysqlInfo.Salt),
		AdminMysqlManager: mysql.NewAdminManager(db, conf.GlobalServerConfig.MysqlInfo.Salt),
		BlobManager:       client.GetBlobService(),
		TokenGenerator:    tg,
	}
	return userservice.NewServer(uss,
		server.WithServiceAddr(utils.NewNetAddr(consts.TCP, net.JoinHostPort(ip, strconv.Itoa(port)))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: conf.GlobalServerConfig.Name}),
	)
}

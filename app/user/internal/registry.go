package internal

import (
	"net"
	"strconv"

	"github.com/alimy/freecar/app/user/conf"
	"github.com/alimy/freecar/library/core/consts"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/hashicorp/consul/api"
	consul "github.com/kitex-contrib/registry-consul"
)

// InitRegistry to init consul
func InitRegistry(Port int) (registry.Registry, *registry.Info) {
	r, err := consul.NewConsulRegister(net.JoinHostPort(
		conf.GlobalConsulConfig.Host,
		strconv.Itoa(conf.GlobalConsulConfig.Port)),
		consul.WithCheck(&api.AgentServiceCheck{
			Interval:                       consts.ConsulCheckInterval,
			Timeout:                        consts.ConsulCheckTimeout,
			DeregisterCriticalServiceAfter: consts.ConsulCheckDeregisterCriticalServiceAfter,
		}))
	if err != nil {
		klog.Fatalf("new consul register failed: %s", err.Error())
	}

	// Using snowflake to generate service name.
	sf, err := snowflake.NewNode(2)
	if err != nil {
		klog.Fatalf("generate service name failed: %s", err.Error())
	}
	info := &registry.Info{
		ServiceName: conf.GlobalServerConfig.Name,
		Addr:        utils.NewNetAddr(consts.TCP, net.JoinHostPort(conf.GlobalServerConfig.Host, strconv.Itoa(Port))),
		Tags: map[string]string{
			"ID": sf.Generate().Base36(),
		},
	}
	return r, info
}

package internal

import (
	"net"
	"strconv"

	"github.com/alimy/freecar/app/api/conf"
	"github.com/alimy/freecar/library/core/consts"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hashicorp/consul/api"
	"github.com/hertz-contrib/registry/consul"
)

// initRegistry to init consul
func initRegistry() (registry.Registry, *registry.Info) {
	// build a consul client
	cfg := api.DefaultConfig()
	cfg.Address = net.JoinHostPort(
		conf.GlobalConsulConfig.Host,
		strconv.Itoa(conf.GlobalConsulConfig.Port))
	consulClient, err := api.NewClient(cfg)
	if err != nil {
		hlog.Fatalf("new consul client failed: %s", err.Error())
	}

	r := consul.NewConsulRegister(consulClient,
		consul.WithCheck(&api.AgentServiceCheck{
			Interval:                       consts.ConsulCheckInterval,
			Timeout:                        consts.ConsulCheckTimeout,
			DeregisterCriticalServiceAfter: consts.ConsulCheckDeregisterCriticalServiceAfter,
		}))

	// Using snowflake to generate service name.
	sf, err := snowflake.NewNode(2)
	if err != nil {
		hlog.Fatalf("generate service name failed: %s", err.Error())
	}
	info := &registry.Info{
		ServiceName: conf.GlobalServerConfig.Name,
		Addr: utils.NewNetAddr(consts.TCP, net.JoinHostPort(conf.GlobalServerConfig.Host,
			strconv.Itoa(conf.GlobalServerConfig.Port))),
		Tags: map[string]string{
			"ID": sf.Generate().Base36(),
		},
		Weight: registry.DefaultWeight,
	}
	return r, info
}

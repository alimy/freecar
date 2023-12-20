package internal

import (
	"net"
	"strconv"

	"github.com/alimy/freecar/app/blob/config"
	"github.com/alimy/freecar/library/core/consts"
	"github.com/alimy/freecar/library/core/utils"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

// initConfig to init consul config server
func initConfig() {
	v := viper.New()
	v.SetConfigFile(consts.BlobConfigPath)
	if err := v.ReadInConfig(); err != nil {
		hlog.Fatalf("read viper config failed: %s", err.Error())
	}
	if err := v.Unmarshal(&config.GlobalConsulConfig); err != nil {
		hlog.Fatalf("unmarshal err failed: %s", err.Error())
	}
	hlog.Infof("Config Info: %v", config.GlobalConsulConfig)

	cfg := api.DefaultConfig()
	cfg.Address = net.JoinHostPort(
		config.GlobalConsulConfig.Host,
		strconv.Itoa(config.GlobalConsulConfig.Port))
	consulClient, err := api.NewClient(cfg)
	if err != nil {
		hlog.Fatalf("new consul client failed: %s", err.Error())
	}
	content, _, err := consulClient.KV().Get(config.GlobalConsulConfig.Key, nil)
	if err != nil {
		hlog.Fatalf("consul kv failed: %s", err.Error())
	}

	err = sonic.Unmarshal(content.Value, &config.GlobalServerConfig)
	if err != nil {
		hlog.Fatalf("sonic unmarshal config failed: %s", err.Error())
	}

	if config.GlobalServerConfig.Host == "" {
		config.GlobalServerConfig.Host, err = utils.GetLocalIPv4Address()
		if err != nil {
			hlog.Fatalf("get localIpv4Addr failed:%s", err.Error())
		}
	}
}

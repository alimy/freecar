package conf

import (
	"net"
	"strconv"

	"github.com/alimy/freecar/library/core/consts"
	"github.com/alimy/freecar/library/core/utils"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

func Initial() {
	v := viper.New()
	v.SetConfigFile(consts.BlobConfigPath)
	if err := v.ReadInConfig(); err != nil {
		klog.Fatalf("read viper config failed: %s", err.Error())
	}
	if err := v.Unmarshal(&GlobalConsulConfig); err != nil {
		klog.Fatalf("unmarshal err failed: %s", err.Error())
	}
	klog.Infof("Config Info: %v", GlobalConsulConfig)

	cfg := api.DefaultConfig()
	cfg.Address = net.JoinHostPort(
		GlobalConsulConfig.Host,
		strconv.Itoa(GlobalConsulConfig.Port))
	consulClient, err := api.NewClient(cfg)
	if err != nil {
		klog.Fatalf("new consul client failed: %s", err.Error())
	}
	content, _, err := consulClient.KV().Get(GlobalConsulConfig.Key, nil)
	if err != nil {
		klog.Fatalf("consul kv failed: %s", err.Error())
	}

	err = sonic.Unmarshal(content.Value, &GlobalServerConfig)
	if err != nil {
		klog.Fatalf("sonic unmarshal config failed: %s", err.Error())
	}

	if GlobalServerConfig.Host == "" {
		GlobalServerConfig.Host, err = utils.GetLocalIPv4Address()
		if err != nil {
			klog.Fatalf("get localIpv4Addr failed:%s", err.Error())
		}
	}
}

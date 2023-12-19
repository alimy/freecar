package config

import (
	"github.com/alimy/freecar/idle/auto/rpc/car/carservice"
	"github.com/alimy/freecar/idle/auto/rpc/profile/profileservice"
	"github.com/alimy/freecar/idle/auto/rpc/user/userservice"
)

var (
	GlobalServerConfig ServerConfig
	GlobalConsulConfig ConsulConfig

	CarClient     carservice.Client
	ProfileClient profileservice.Client
	UserClient    userservice.Client
)

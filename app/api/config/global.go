package config

import (
	"github.com/alimy/freecar/idle/auto/rpc/car/carservice"
	"github.com/alimy/freecar/idle/auto/rpc/profile/profileservice"
	"github.com/alimy/freecar/idle/auto/rpc/trip/tripservice"
	"github.com/alimy/freecar/idle/auto/rpc/user/userservice"
)

var (
	GlobalServerConfig ServerConfig
	GlobalConsulConfig ConsulConfig

	GlobalUserClient    userservice.Client
	GlobalCarClient     carservice.Client
	GlobalProfileClient profileservice.Client
	GlobalTripClient    tripservice.Client
)

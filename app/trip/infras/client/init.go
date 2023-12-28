package client

import (
	"sync"

	"github.com/alimy/freecar/idle/auto/rpc/car/carservice"
	"github.com/alimy/freecar/idle/auto/rpc/profile/profileservice"
	"github.com/alimy/freecar/idle/auto/rpc/user/userservice"
)

var (
	userSvc    userservice.Client
	carSvc     carservice.Client
	profileSvc profileservice.Client
	_once      sync.Once
)

// initial initialize rpc service
func initial() {
	_once.Do(func() {
		userSvc = initUser()
		carSvc = initCar()
		profileSvc = initProfile()
	})
}

func GetUserService() userservice.Client {
	initial()
	return userSvc
}

func GetCarService() carservice.Client {
	initial()
	return carSvc
}

func GetProfileService() profileservice.Client {
	initial()
	return profileSvc
}

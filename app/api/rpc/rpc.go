package rpc

import (
	"github.com/alimy/freecar/idle/auto/rpc/car/carservice"
	"github.com/alimy/freecar/idle/auto/rpc/profile/profileservice"
	"github.com/alimy/freecar/idle/auto/rpc/trip/tripservice"
	"github.com/alimy/freecar/idle/auto/rpc/user/userservice"
)

var (
	UserSvc    userservice.Client
	CarSvc     carservice.Client
	ProfileSvc profileservice.Client
	TripSvc    tripservice.Client
)

// Initial initialize rpc service
func Initial() {
	initUser()
	initCar()
	initProfile()
	initTrip()
}

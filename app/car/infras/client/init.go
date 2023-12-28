package client

import (
	"sync"

	"github.com/alimy/freecar/idle/auto/rpc/car/carservice"
	"github.com/alimy/freecar/idle/auto/rpc/trip/tripservice"
)

var (
	CarSvc  carservice.Client
	TripSvc tripservice.Client
	_once   sync.Once
)

// Initial initialize rpc service
func Initial() {
	_once.Do(func() {
		initCar()
		initTrip()
	})
}

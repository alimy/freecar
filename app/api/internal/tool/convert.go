package tool

import (
	hbase "github.com/alimy/freecar/app/api/biz/model/base"
	kbase "github.com/alimy/freecar/idle/auto/rpc/base"
)

func ConvertTripLocation(l *hbase.Location) *kbase.Location {
	return &kbase.Location{
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
	}
}

func ConvertDriver(d *hbase.Driver) *kbase.Driver {
	if d == nil {
		return nil
	}
	return &kbase.Driver{
		Id:        d.ID,
		AvatarUrl: d.AvatarURL,
	}
}

func ConvertCarLocation(l *hbase.Position) *kbase.Position {
	if l == nil {
		return nil
	}
	return &kbase.Position{
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
	}
}

func ConvertCar(c *hbase.Car) *kbase.Car {
	if c == nil {
		return nil
	}
	return &kbase.Car{
		Status:   kbase.CarStatus(c.Status),
		Driver:   ConvertDriver(c.Driver),
		Position: ConvertCarLocation(c.Position),
		TripId:   c.TripID,
		Power:    c.Power,
		PlateNum: c.PlateNum,
	}
}

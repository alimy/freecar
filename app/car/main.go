package main

import (
	"context"

	"github.com/alimy/freecar/app/car/config"
	"github.com/alimy/freecar/app/car/internal"
	"github.com/alimy/freecar/app/car/pkg/mq/amqpclt"
	"github.com/alimy/freecar/app/car/pkg/sim"
	"github.com/alimy/freecar/app/car/pkg/trip"
	"github.com/alimy/freecar/app/car/pkg/ws"
	"github.com/alimy/freecar/app/car/rpc"
	"github.com/alimy/freecar/app/car/servants"
	hzserver "github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
)

func main() {
	// initialization
	internal.Initial()
	rpc.Initial()
	amqpC := internal.InitMq()
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	mqInfo := config.GlobalServerConfig.RabbitMqInfo
	subscriber, err := amqpclt.NewSubscriber(amqpC, mqInfo.Exchange)
	if err != nil {
		klog.Fatal("cannot create subscriber")
	}

	// Create new server.
	srv := servants.NewCarService()

	h := hzserver.Default(hzserver.WithHostPorts(config.GlobalServerConfig.WsAddr))
	h.GET("/ws", ws.Handler(subscriber))
	h.NoHijackConnPool = true
	go func() {
		klog.Infof("HTTP server started. addr: %s", config.GlobalServerConfig.WsAddr)
		h.Spin()
	}()

	go trip.RunUpdater(subscriber, rpc.TripSvc)

	simController := sim.Controller{
		CarService: rpc.CarSvc,
		Subscriber: subscriber,
	}
	go simController.RunSimulations(context.Background())

	err = srv.Run()
	if err != nil {
		klog.Fatal(err)
	}
}

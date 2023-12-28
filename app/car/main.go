package main

import (
	"context"

	"github.com/alimy/freecar/app/car/conf"
	"github.com/alimy/freecar/app/car/infras/client"
	"github.com/alimy/freecar/app/car/infras/mq/amqpclt"
	"github.com/alimy/freecar/app/car/infras/sim"
	"github.com/alimy/freecar/app/car/infras/trip"
	"github.com/alimy/freecar/app/car/infras/ws"
	"github.com/alimy/freecar/app/car/internal"
	"github.com/alimy/freecar/app/car/servants"
	hzserver "github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
)

func main() {
	// initialization
	conf.Initial()
	internal.Initial()
	client.Initial()
	amqpC := internal.InitMq()
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(conf.GlobalServerConfig.Name),
		provider.WithExportEndpoint(conf.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	mqInfo := conf.GlobalServerConfig.RabbitMqInfo
	subscriber, err := amqpclt.NewSubscriber(amqpC, mqInfo.Exchange)
	if err != nil {
		klog.Fatal("cannot create subscriber")
	}

	// Create new server.
	srv := servants.NewCarService()

	h := hzserver.Default(hzserver.WithHostPorts(conf.GlobalServerConfig.WsAddr))
	h.GET("/ws", ws.Handler(subscriber))
	h.NoHijackConnPool = true
	go func() {
		klog.Infof("HTTP server started. addr: %s", conf.GlobalServerConfig.WsAddr)
		h.Spin()
	}()

	go trip.RunUpdater(subscriber, client.TripSvc)

	simController := sim.Controller{
		CarService: client.CarSvc,
		Subscriber: subscriber,
	}
	go simController.RunSimulations(context.Background())

	err = srv.Run()
	if err != nil {
		klog.Fatal(err)
	}
}

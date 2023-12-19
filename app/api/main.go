// Code generated by hertz generator.

package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/alimy/freecar/app/api/biz/handler"
	"github.com/alimy/freecar/app/api/biz/router"
	"github.com/alimy/freecar/app/api/config"
	"github.com/alimy/freecar/app/api/internal"
	"github.com/alimy/freecar/app/api/rpc"
	"github.com/alimy/freecar/library/core/errno"
	"github.com/alimy/freecar/library/core/tools"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	cfg "github.com/hertz-contrib/http2/config"
	"github.com/hertz-contrib/http2/factory"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	hertzSentinel "github.com/hertz-contrib/opensergo/sentinel/adapter"
	"github.com/hertz-contrib/pprof"
)

func main() {
	// initialize
	r, info, tlsCfg := internal.Initial()
	rpc.Initial()
	tracer, trcCfg := hertztracing.NewServerTracer()
	// create a new server
	h := server.New(
		tracer,
		server.WithALPN(true),
		server.WithTLS(tlsCfg),
		server.WithHostPorts(fmt.Sprintf(":%d", config.GlobalServerConfig.Port)),
		server.WithRegistry(r, info),
		server.WithHandleMethodNotAllowed(true),
	)
	// add h2
	h.AddProtocol("h2", factory.NewServerFactory(
		cfg.WithReadTimeout(time.Minute),
		cfg.WithDisableKeepAlive(false)))
	tlsCfg.NextProtos = append(tlsCfg.NextProtos, "h2")
	// use pprof & tracer & sentinel
	pprof.Register(h)
	h.Use(hertztracing.ServerMiddleware(trcCfg))
	h.Use(hertzSentinel.SentinelServerMiddleware(
		// abort with status 429 by default
		hertzSentinel.WithServerBlockFallback(func(c context.Context, ctx *app.RequestContext) {
			ctx.JSON(http.StatusTooManyRequests, nil)
			ctx.Abort()
		}),
	))
	register(h)
	h.Spin()
}

// register registers all routers.
func register(r *server.Hertz) {
	router.GeneratedRegister(r)
	customizedRegister(r)
}

// customizeRegister registers customize routers.
func customizedRegister(r *server.Hertz) {
	r.GET("/ping", handler.Ping)

	// your code ...
	r.NoRoute(func(ctx context.Context, c *app.RequestContext) { // used for HTTP 404
		c.JSON(http.StatusNotFound, tools.BuildBaseResp(errno.NoRoute))
	})
	r.NoMethod(func(ctx context.Context, c *app.RequestContext) { // used for HTTP 405
		c.JSON(http.StatusMethodNotAllowed, tools.BuildBaseResp(errno.NoMethod))
	})
}

package internal

import (
	"crypto/tls"

	"github.com/alimy/freecar/app/api/config"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
)

func Initial() (r registry.Registry, i *registry.Info, c *tls.Config) {
	initLogger()
	config.Initial()
	r, i = initRegistry()
	initSentinel()
	c = initTLS()
	return
}

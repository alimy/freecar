package internal

import (
	"crypto/tls"

	"github.com/cloudwego/hertz/pkg/app/server/registry"
)

func Initial() (r registry.Registry, i *registry.Info, c *tls.Config) {
	initLogger()
	initConfig()
	r, i = initRegistry()
	initSentinel()
	c = initTLS()
	return
}

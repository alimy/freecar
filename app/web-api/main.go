package main

import (
	"github.com/alimy/freecar/app/web-api/servants"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	h := server.Default()

	// register servants to hertz
	servants.RegisterServants(h.Engine)

	// start servant service
	h.Spin()
}

package main

import (
	"log"

	"github.com/alimy/freecar/app/admin/servants"
	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.Default()

	// register servants to gin
	servants.RegisterServants(e)

	// start servant service
	if err := e.Run(); err != nil {
		log.Fatal(err)
	}
}

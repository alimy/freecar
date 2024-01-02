//go:build generate
// +build generate

package main

import (
	"context"
	"log"

	. "github.com/alimy/mir/v4/core"
	. "github.com/alimy/mir/v4/engine"
	"github.com/cloudwego/hertz/pkg/app"

	_ "github.com/alimy/freecar/api/web-api/mirc/admin/v1"
	_ "github.com/alimy/freecar/api/web-api/mirc/web/v1"
	_ "github.com/alimy/freecar/api/web-api/mirc/web/v2"
)

//go:generate go run $GOFILE
func main() {
	log.Println("[Mir] generate code start")
	opts := Options{
		UseHertz(),
		SinkPath("../auto"),
		WatchCtxDone(true),
		RunMode(InSerialMode),
		AssertType4[context.Context, *app.RequestContext](),
	}
	if err := Generate(opts); err != nil {
		log.Fatal(err)
	}
	log.Println("[Mir] generate code finish")
}

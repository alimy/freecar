package mq

import (
	"context"

	car "github.com/alimy/freecar/idle/auto/rpc/base"
)

// Publisher defines the publish interface.
type Publisher interface {
	Publish(context.Context, *car.CarEntity) error
}

// Subscriber defines a car update subscriber.
type Subscriber interface {
	Subscribe(context.Context) (ch chan *car.CarEntity, cleanUp func(), err error)
}

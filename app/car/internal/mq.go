package internal

import (
	"fmt"

	"github.com/alimy/freecar/app/car/conf"
	"github.com/alimy/freecar/library/core/consts"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/streadway/amqp"
)

// InitMq to init rabbitMQ
func InitMq() *amqp.Connection {
	c := conf.GlobalServerConfig.RabbitMqInfo
	amqpConn, err := amqp.Dial(fmt.Sprintf(consts.RabbitMqURI, c.User, c.Password, c.Host, c.Port))
	if err != nil {
		klog.Fatal("cannot dial amqp", err)
	}
	return amqpConn
}

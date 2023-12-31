package internal

import (
	"fmt"
	"time"

	"github.com/alimy/freecar/app/blob/conf"
	"github.com/alimy/freecar/library/core/consts"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/opentelemetry/logging/logrus"
	"gorm.io/plugin/opentelemetry/tracing"
)

func NewDB() *gorm.DB {
	c := conf.GlobalServerConfig.MysqlInfo
	dsn := fmt.Sprintf(consts.MySqlDSN, c.User, c.Password, c.Host, c.Port, c.Name)
	newLogger := logger.New(
		logrus.NewWriter(), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL Threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,          // Disable color printing
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		klog.Fatalf("init gorm failed: %s", err.Error())
	}
	if err := db.Use(tracing.NewPlugin()); err != nil {
		klog.Fatalf("use tracing plugin failed: %s", err.Error())
	}
	return db
}

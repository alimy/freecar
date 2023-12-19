package internal

import (
	"crypto/tls"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// initTLS
func initTLS() *tls.Config {
	cfg := &tls.Config{
		MinVersion:         tls.VersionTLS10,
		InsecureSkipVerify: true,
	}
	cert, err := tls.LoadX509KeyPair("cert/server.crt",
		"server.key")
	if err != nil {
		hlog.Fatal("tls failed", err)
	}
	cfg.Certificates = append(cfg.Certificates, cert)
	return cfg
}

package rpc

import (
	"sync"

	"github.com/alimy/freecar/idle/auto/rpc/blob/blobservice"
)

var (
	blobSvc blobservice.Client
	_once   sync.Once
)

// initial initialize rpc service
func initial() {
	_once.Do(func() {
		blobSvc = initBlob()
	})
}

func GetBlobService() blobservice.Client {
	initial()
	return blobSvc
}

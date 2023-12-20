package rpc

import (
	"sync"

	"github.com/alimy/freecar/idle/auto/rpc/blob/blobservice"
)

var (
	BlobSvc blobservice.Client
	_once   sync.Once
)

// Initial initialize rpc service
func Initial() {
	_once.Do(func() {
		initBlob()
	})
}

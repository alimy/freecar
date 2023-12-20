package servants

import (
	"context"
	"fmt"
	"time"

	"github.com/alimy/freecar/app/blob/pkg/minio"
	"github.com/alimy/freecar/app/blob/pkg/mysql"
	"github.com/alimy/freecar/app/blob/pkg/redis"
	"github.com/alimy/freecar/idle/auto/rpc/blob"
	"github.com/alimy/freecar/library/core/consts"
	"github.com/alimy/freecar/library/core/errno"
	"github.com/alimy/freecar/library/core/id"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
)

// blobSrv implements the last service interface defined in the IDL.
type blobSrv struct {
	minioManager *minio.Manager
	mysqlManager *mysql.Manager
	redisManager *redis.Manager
}

// CreateBlob implements the blobSrv interface.
func (s *blobSrv) CreateBlob(ctx context.Context, req *blob.CreateBlobRequest) (*blob.CreateBlobResponse, error) {
	var br mysql.BlobRecord
	br.AccountId = req.AccountId

	sf, err := snowflake.NewNode(consts.BlobSnowflakeNode)
	if err != nil {
		klog.Fatalf("generate id failed: %s", err.Error())
	}
	br.Path = fmt.Sprintf("%s/%s", req.AccountId, sf.Generate().String())

	err = s.mysqlManager.CreateBlobRecord(&br)
	if err != nil {
		klog.Error("create blob record err", err)
		return nil, errno.BlobSrvErr
	}
	url, err := s.minioManager.PutObjectURL(ctx, br.Path, time.Duration(req.UploadUrlTimeoutSec)*time.Second)
	if err != nil {
		klog.Error("presigned put object url err", err)
		return nil, errno.BlobSrvErr
	}
	return &blob.CreateBlobResponse{
		Id:        br.ID,
		UploadUrl: url,
	}, nil
}

// GetBlobURL implements the blobSrv interface.
func (s *blobSrv) GetBlobURL(ctx context.Context, req *blob.GetBlobURLRequest) (*blob.GetBlobURLResponse, error) {
	br, err := s.redisManager.Get(ctx, id.BlobID(req.Id))
	if err != nil {
		klog.Error("get blob cache err", err)
		br, err = s.mysqlManager.GetBlobRecord(req.Id)
		if err == errno.RecordNotFound {
			return nil, errno.RecordNotFound
		}
		if err != nil {
			klog.Error("get blob record err", err)
			return nil, errno.BlobSrvErr.WithMessage("get blob record err")
		}
		go func() {
			if err := s.redisManager.Insert(context.Background(), br); err != nil {
				klog.Error("create cache record err", err)
			}
		}()
	}
	url, err := s.minioManager.GetObjectURL(ctx, br.Path, time.Duration(req.TimeoutSec)*time.Second)
	if err != nil {
		klog.Error("cannot get object url", err)
		return nil, errno.BlobSrvErr
	}

	return &blob.GetBlobURLResponse{Url: url}, nil
}

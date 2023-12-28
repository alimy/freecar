package servants

import (
	"hash/fnv"

	"github.com/alimy/freecar/idle/auto/rpc/base"
	"github.com/bytedance/sonic"
)

var _poi = []string{
	"知行苑7舍",
	"兴业苑5舍",
	"中心食堂",
	"第二教学楼",
	"综合实验大楼",
	"信科大厦",
}

type poiManager struct{}

// Resolve resolves the given location.
func (*poiManager) Resolve(loc *base.Location) (string, error) {
	b, err := sonic.Marshal(loc)
	if err != nil {
		return "", err
	}
	h := fnv.New32()
	h.Write(b)
	return _poi[int(h.Sum32())%len(_poi)], nil
}

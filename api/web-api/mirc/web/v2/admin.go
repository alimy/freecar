package v1

import (
	"github.com/alimy/freecar/idle/model/web"
	. "github.com/alimy/mir/v4"
	. "github.com/alimy/mir/v4/engine"
)

func init() {
	Entry[Admin]()
}

// Admin 运维相关服务
type Admin struct {
	Chain `mir:"-"`
	Group `mir:"v2"`

	// ChangeUserStatus 管理·禁言/解封用户
	ChangeUserStatus func(Post, web.ChangeUserStatusReq)                  `mir:"/admin/user/status"`
	SiteInfo         func(Get, Context, web.SiteInfoReq) web.SiteInfoResp `mir:"/admin/site/status"`
}

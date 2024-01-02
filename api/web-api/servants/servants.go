package servants

import (
	api "github.com/alimy/freecar/api/web-api/auto/api/v1"
	apiv2 "github.com/alimy/freecar/api/web-api/auto/api/v2"
	"github.com/cloudwego/hertz/pkg/route"
)

// RegisterServants register all the servants to gin.Engine
func RegisterServants(e *route.Engine) {
	api.RegisterSiteServant(e, newSiteSrvA())
	api.RegisterAdminServant(e, newAdminSrvA())
	apiv2.RegisterAdminServant(e, newAdminSrvB())
	apiv2.RegisterSiteServant(e, newSiteSrvB())
}

package servants

import (
	api "github.com/alimy/freecar/api/admin/auto/api/v1"
	apiv2 "github.com/alimy/freecar/api/admin/auto/api/v2"
	"github.com/gin-gonic/gin"
)

// RegisterServants register all the servants to gin.Engine
func RegisterServants(e *gin.Engine) {
	api.RegisterSiteServant(e, newSiteSrvA())
	api.RegisterAdminServant(e, newAdminSrvA())
	apiv2.RegisterAdminServant(e, newAdminSrvB())
	apiv2.RegisterSiteServant(e, newSiteSrvB())
}

package servants

import (
	api "github.com/alimy/freecar/app/admin/auto/api/v1"
	"github.com/gin-gonic/gin"
)

// RegisterServants register all the servants to gin.Engine
func RegisterServants(e *gin.Engine) {
	api.RegisterSiteServant(e, newSiteSrvA())
	api.RegisterAdminServant(e, newAdminSrvA())
}

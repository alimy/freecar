package servants

import (
	api "github.com/alimy/freecar/api/admin/auto/api/v1"
)

type siteSrvA struct {
	baseSrv

	api.UnimplementedSiteServant
}

func newSiteSrvA() api.Site {
	return &siteSrvA{}
}

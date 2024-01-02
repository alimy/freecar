package servants

import (
	api "github.com/alimy/freecar/api/web-api/auto/api/v2"
)

type siteSrvB struct {
	baseSrv

	api.UnimplementedSiteServant
}

func newSiteSrvB() api.Site {
	return &siteSrvB{}
}

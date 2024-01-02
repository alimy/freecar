package servants

import (
	api "github.com/alimy/freecar/api/admin/auto/api/v2"
)

type adminSrvB struct {
	baseSrv

	api.UnimplementedAdminServant
}

func newAdminSrvB() api.Admin {
	return &adminSrvB{}
}

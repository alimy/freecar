package servants

import (
	api "github.com/alimy/freecar/api/admin/auto/api/v1"
)

type adminSrvA struct {
	baseSrv

	api.UnimplementedAdminServant
}

func newAdminSrvA() api.Admin {
	return &adminSrvA{}
}

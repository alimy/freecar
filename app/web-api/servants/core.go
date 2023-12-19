package servants

import (
	"context"
	"net/http"

	"github.com/alimy/mir/v4"
	"github.com/cloudwego/hertz/pkg/app"
)

type baseSrv struct{}

func (baseSrv) Bind(c context.Context, ctx *app.RequestContext, obj any) mir.Error {
	if err := ctx.BindAndValidate(obj); err != nil {
		mir.NewError(http.StatusBadRequest, err)
	}
	return nil
}

func (baseSrv) Render(c context.Context, ctx *app.RequestContext, data any, err mir.Error) {
	if err == nil {
		ctx.JSON(http.StatusOK, data)
	} else {
		ctx.JSON(err.StatusCode(), err.Error())
	}
	ctx.String(http.StatusNotImplemented, "")
}

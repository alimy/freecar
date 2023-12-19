// Code generated by go-mir. DO NOT EDIT.
// versions:
// - mir v4.1.0

package v2

import (
	"context"
	"net/http"

	"github.com/alimy/freecar/idle/model/web"
	"github.com/alimy/mir/v4"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
)

type Site interface {
	_default_

	Profile() (*web.SiteProfileResp, mir.Error)
	Version() (*web.VersionResp, mir.Error)

	mustEmbedUnimplementedSiteServant()
}

// RegisterSiteServant register Site servant to hertz
func RegisterSiteServant(e *route.Engine, s Site) {
	router := e.Group("v2")

	// register routes info to router
	router.Handle("GET", "/site/profile", func(c context.Context, ctx *app.RequestContext) {
		select {
		case <-c.Done():
			return
		default:
		}

		resp, err := s.Profile()
		s.Render(c, ctx, resp, err)
	})
	router.Handle("GET", "/site/version", func(c context.Context, ctx *app.RequestContext) {
		select {
		case <-c.Done():
			return
		default:
		}

		resp, err := s.Version()
		s.Render(c, ctx, resp, err)
	})
}

// UnimplementedSiteServant can be embedded to have forward compatible implementations.
type UnimplementedSiteServant struct{}

func (UnimplementedSiteServant) Profile() (*web.SiteProfileResp, mir.Error) {
	return nil, mir.Errorln(http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
}

func (UnimplementedSiteServant) Version() (*web.VersionResp, mir.Error) {
	return nil, mir.Errorln(http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
}

func (UnimplementedSiteServant) mustEmbedUnimplementedSiteServant() {}

// Code generated by go-mir. DO NOT EDIT.
// versions:
// - mir v4.1.0

package v2

import (
	"net/http"

	"github.com/alimy/freecar/idle/model/web"
	"github.com/alimy/mir/v4"
	"github.com/gin-gonic/gin"
)

type Site interface {
	_default_

	Profile() (*web.SiteProfileResp, mir.Error)
	Version() (*web.VersionResp, mir.Error)

	mustEmbedUnimplementedSiteServant()
}

// RegisterSiteServant register Site servant to gin
func RegisterSiteServant(e *gin.Engine, s Site) {
	router := e.Group("v2")

	// register routes info to router
	router.Handle("GET", "/site/profile", func(c *gin.Context) {
		select {
		case <-c.Request.Context().Done():
			return
		default:
		}

		resp, err := s.Profile()
		s.Render(c, resp, err)
	})
	router.Handle("GET", "/site/version", func(c *gin.Context) {
		select {
		case <-c.Request.Context().Done():
			return
		default:
		}

		resp, err := s.Version()
		s.Render(c, resp, err)
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

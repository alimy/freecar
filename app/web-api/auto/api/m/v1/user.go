// Code generated by go-mir. DO NOT EDIT.
// versions:
// - mir v4.1.0

package v1

import (
	"context"
	"net/http"

	"github.com/alimy/mir/v4"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
)

type _binding_ interface {
	Bind(context.Context, *app.RequestContext) mir.Error
}

type _render_ interface {
	Render(context.Context, *app.RequestContext)
}

type _default_ interface {
	Bind(context.Context, *app.RequestContext, any) mir.Error
	Render(context.Context, *app.RequestContext, any, mir.Error)
}

type LoginReq struct {
	AgentInfo AgentInfo `json:"agent_info"`
	Name      string    `json:"name"`
	Passwd    string    `json:"passwd"`
}

type AgentInfo struct {
	Platform  string `json:"platform"`
	UserAgent string `json:"user_agent"`
}

type LoginResp struct {
	UserInfo
	ServerInfo ServerInfo `json:"server_info"`
	JwtToken   string     `json:"jwt_token"`
}

type ServerInfo struct {
	ApiVer string `json:"api_ver"`
}

type UserInfo struct {
	Name string `json:"name"`
}

type User interface {
	_default_

	// Chain provide handlers chain for hertz
	Chain() []app.HandlerFunc

	Logout() mir.Error
	Login(*LoginReq) (*LoginResp, mir.Error)

	mustEmbedUnimplementedUserServant()
}

// RegisterUserServant register User servant to hertz
func RegisterUserServant(e *route.Engine, s User) {
	router := e.Group("m/v1")
	// use chain for router
	middlewares := s.Chain()
	router.Use(middlewares...)

	// register routes info to router
	router.Handle("POST", "/user/logout/", func(c context.Context, ctx *app.RequestContext) {
		select {
		case <-c.Done():
			return
		default:
		}

		s.Render(c, ctx, nil, s.Logout())
	})
	router.Handle("POST", "/user/login/", func(c context.Context, ctx *app.RequestContext) {
		select {
		case <-c.Done():
			return
		default:
		}
		req := new(LoginReq)
		if err := s.Bind(c, ctx, req); err != nil {
			s.Render(c, ctx, nil, err)
			return
		}
		resp, err := s.Login(req)
		s.Render(c, ctx, resp, err)
	})
}

// UnimplementedUserServant can be embedded to have forward compatible implementations.
type UnimplementedUserServant struct{}

func (UnimplementedUserServant) Chain() []app.HandlerFunc {
	return nil
}

func (UnimplementedUserServant) Logout() mir.Error {
	return mir.Errorln(http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
}

func (UnimplementedUserServant) Login(req *LoginReq) (*LoginResp, mir.Error) {
	return nil, mir.Errorln(http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
}

func (UnimplementedUserServant) mustEmbedUnimplementedUserServant() {}

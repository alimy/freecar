package common

import (
	"context"
	"net/http"

	pt "aidanwoods.dev/go-paseto"
	"github.com/alimy/freecar/app/api/config"
	"github.com/alimy/freecar/library/core/consts"
	"github.com/alimy/freecar/library/core/errno"
	"github.com/alimy/freecar/library/core/middleware"
	"github.com/alimy/freecar/library/core/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/gzip"
	"github.com/hertz-contrib/paseto"
)

func CommonMW() []app.HandlerFunc {
	return []app.HandlerFunc{
		// use cors mw
		middleware.Cors(),
		// use recovery mw
		middleware.Recovery(),
		// use gzip mw
		gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedExtensions([]string{".jpg", ".mp4", ".png"})),
	}
}

func PasetoAuth(audience string) app.HandlerFunc {
	pi := config.GlobalServerConfig.PasetoInfo
	pf, err := paseto.NewV4PublicParseFunc(pi.PubKey, []byte(pi.Implicit), paseto.WithAudience(audience), paseto.WithNotBefore())
	if err != nil {
		hlog.Fatal(err)
	}
	sh := func(ctx context.Context, c *app.RequestContext, token *pt.Token) {
		aid, err := token.GetString("id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.BuildBaseResp(errno.BadRequest.WithMessage("missing accountID in token")))
			c.Abort()
			return
		}
		c.Set(consts.AccountID, aid)
	}

	eh := func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusUnauthorized, utils.BuildBaseResp(errno.BadRequest.WithMessage("invalid token")))
		c.Abort()
	}
	return paseto.New(paseto.WithTokenPrefix("Bearer "), paseto.WithParseFunc(pf), paseto.WithSuccessHandler(sh), paseto.WithErrorFunc(eh))
}

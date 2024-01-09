// Code generated by hertz generator.

package gpt

import (
	"context"
	"net/http"

	"github.com/alimy/freecar/api/api/biz/model/base"
	"github.com/alimy/freecar/api/api/biz/model/gpt"
	"github.com/alimy/freecar/api/api/conf"
	sConst "github.com/alimy/freecar/library/core/consts"
	"github.com/alimy/freecar/library/core/errno"
	"github.com/alimy/freecar/library/core/utils"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type requestRaw struct {
	Model       string     `json:"model"`
	Messages    []messages `json:"messages"`
	Temperature float64    `json:"temperature"`
}
type messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type requestRes struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}

// Chat .
// @router /chat [POST]
func Chat(ctx context.Context, c *app.RequestContext) {
	var err error
	var req gpt.ChatRequest
	resp := new(gpt.ChatResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = (*base.BaseResponse)(utils.BuildBaseResp(errno.ParamsErr))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	clt, err := client.NewClient(client.WithDialer(standard.NewDialer()))
	if err != nil {
		hlog.Error("new hertz client error", err)
		resp.BaseResp = (*base.BaseResponse)(utils.BuildBaseResp(errno.BadRequest))
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	clt.SetProxy(protocol.ProxyURI(protocol.ParseURI(conf.GlobalServerConfig.ProxyURL)))

	reqRaw := &requestRaw{
		Model: "gpt-3.5-turbo",
		Messages: []messages{
			{
				Role:    "user",
				Content: req.Content,
			},
		},
		Temperature: 0.7,
	}
	hReq := &protocol.Request{}
	hRes := &protocol.Response{}
	hReq.SetMethod(consts.MethodPost)
	hReq.Header.SetContentTypeBytes([]byte("application/json"))
	hReq.SetRequestURI(sConst.GPTUrl)
	hReq.SetHeader("Authorization", "Bearer "+conf.GlobalServerConfig.GPTKey)
	body, err := sonic.Marshal(reqRaw)
	hReq.SetBody(body)
	err = clt.Do(ctx, hReq, hRes)
	if err != nil {
		hlog.Error("chat with gpt error", err)
		resp.BaseResp = (*base.BaseResponse)(utils.BuildBaseResp(errno.BadRequest))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	var res requestRes
	err = sonic.Unmarshal(hRes.Body(), &res)
	if err != nil {
		hlog.Error("unmarshal error", err)
		resp.BaseResp = (*base.BaseResponse)(utils.BuildBaseResp(errno.BadRequest))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.BaseResp = (*base.BaseResponse)(utils.BuildBaseResp(errno.Success))
	resp.Content = res.Choices[0].Message.Content
	c.JSON(consts.StatusOK, resp)
}
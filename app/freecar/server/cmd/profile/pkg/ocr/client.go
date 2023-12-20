package ocr

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/alimy/freecar/app/api/cmd/profile/config"
	"github.com/alimy/freecar/idle/auto/rpc/base"
	"github.com/alimy/freecar/library/core/consts"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/network/standard"
)

type LicenseManager struct{}

type RequestRes struct {
	WordsResult struct {
		Name struct {
			Words string `json:"words"`
		} `json:"姓名"`
		BirthDay struct {
			Words string `json:"words"`
		} `json:"出生日期"`
		LicenseNum struct {
			Words string `json:"words"`
		} `json:"证号"`
		Gender struct {
			Words string `json:"words"`
		} `json:"性别"`
	} `json:"words_result"`
}

func (l *LicenseManager) GetLicenseInfo(url string) (*base.Identity, error) {
	c, err := client.NewClient(client.WithDialer(standard.NewDialer()))
	if err != nil {
		hlog.Error("new hertz client error", err)
		return nil, err
	}
	ocrUrl := fmt.Sprintf("%s?access_token=%s&url=%s",
		consts.OCRUrl, config.GlobalServerConfig.OCRConfig.AccessToken, url)
	_, body, err := c.Post(context.Background(), nil, ocrUrl, nil)
	if err != nil {
		hlog.Error("get license info error", err)
		return nil, err
	}
	var res RequestRes
	err = sonic.Unmarshal(body, &res)
	if err != nil {
		hlog.Error("unmarshal license info error", err)
		return nil, err
	}
	var gender base.Gender
	if res.WordsResult.Gender.Words == "男" {
		gender = 1
	} else if res.WordsResult.Gender.Words == "女" {
		gender = 2
	} else {
		gender = 0
	}
	sBirth := res.WordsResult.BirthDay.Words
	year, _ := strconv.Atoi(sBirth[0:4])
	month, _ := strconv.Atoi(sBirth[4:6])
	day, _ := strconv.Atoi(sBirth[6:])
	birthDateMillis := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Unix()
	identity := &base.Identity{
		LicNumber:       res.WordsResult.LicenseNum.Words,
		Name:            res.WordsResult.Name.Words,
		Gender:          gender,
		BirthDateMillis: birthDateMillis,
	}
	return identity, nil
}

// Code generated by hertz generator.

package profile

import (
	"context"
	"net/http"

	hprofile "github.com/alimy/freecar/api/api/biz/model/profile"
	"github.com/alimy/freecar/api/api/rpc"
	kbase "github.com/alimy/freecar/idle/auto/rpc/base"
	kprofile "github.com/alimy/freecar/idle/auto/rpc/profile"
	"github.com/alimy/freecar/library/core/consts"
	"github.com/alimy/freecar/library/core/errno"
	"github.com/alimy/freecar/library/core/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// DeleteProfile .
// @router /profile/admin/delete [DELETE]
func DeleteProfile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hprofile.DeleteProfileRequest
	resp := new(kprofile.DeleteProfileResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamsErr)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	res, err := rpc.ProfileSvc.DeleteProfile(ctx, &kprofile.DeleteProfileRequest{AccountId: req.AccountID})
	if err != nil {
		hlog.Error("rpc profile service err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, res)
}

// CheckProfile .
// @router /profile/admin/check [POST]
func CheckProfile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hprofile.CheckProfileRequest
	resp := new(kprofile.CheckProfileResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamsErr)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	res, err := rpc.ProfileSvc.CheckProfile(ctx, &kprofile.CheckProfileRequest{
		AccountId: req.AccountID,
		Accept:    req.Accept,
	})
	if err != nil {
		hlog.Error("rpc profile service err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetAllProfile .
// @router /profile/admin/all [GET]
func GetAllProfile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hprofile.GetAllProfileRequest
	resp := new(kprofile.GetAllProfileResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamsErr)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	res, err := rpc.ProfileSvc.GetAllProfile(ctx, &kprofile.GetAllProfileRequest{})
	if err != nil {
		hlog.Error("rpc profile service err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetSomeProfile .
// @router /profile/admin/some [GET]
func GetSomeProfile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hprofile.GetSomeProfileRequest
	resp := new(kprofile.GetSomeProfileResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamsErr)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	res, err := rpc.ProfileSvc.GetSomeProfile(ctx, &kprofile.GetSomeProfileRequest{})
	if err != nil {
		hlog.Error("rpc profile service err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetPendingProfile .
// @router /profile/admin/pending [GET]
func GetPendingProfile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hprofile.GetPendingProfileRequest
	resp := new(kprofile.GetPendingProfileResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamsErr)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	res, err := rpc.ProfileSvc.GetPendingProfile(ctx, &kprofile.GetPendingProfileRequest{})
	if err != nil {
		hlog.Error("rpc profile service err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetProfile .
// @router /profile/mini/profile [GET]
func GetProfile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hprofile.GetProfileRequest
	resp := new(kprofile.GetProfileResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamsErr)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	res, err := rpc.ProfileSvc.GetProfile(ctx, &kprofile.GetProfileRequest{AccountId: c.MustGet(consts.AccountID).(string)})
	if err != nil {
		hlog.Error("rpc profile service err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, res)
}

// SubmitProfile .
// @router /profile/mini/profile [POST]
func SubmitProfile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hprofile.SubmitProfileRequest
	resp := new(kprofile.SubmitProfileResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamsErr)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	res, err := rpc.ProfileSvc.SubmitProfile(ctx, &kprofile.SubmitProfileRequest{
		AccountId: c.MustGet(consts.AccountID).(string),
		Identity: &kbase.Identity{
			LicNumber:       req.Identity.LicNumber,
			Name:            req.Identity.Name,
			Gender:          kbase.Gender(req.Identity.Gender),
			BirthDateMillis: req.Identity.BirthDateMillis,
		},
	})
	if err != nil {
		hlog.Error("rpc profile service err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, res)
}

// ClearProfile .
// @router /profile/mini/profile [DELETE]
func ClearProfile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hprofile.ClearProfileRequest
	resp := new(kprofile.ClearProfileResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamsErr)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	res, err := rpc.ProfileSvc.ClearProfile(ctx, &kprofile.ClearProfileRequest{AccountId: c.MustGet(consts.AccountID).(string)})
	if err != nil {
		hlog.Error("rpc profile service err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, res)
}

// CreateProfilePhoto .
// @router /profile/mini/photo [POST]
func CreateProfilePhoto(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hprofile.CreateProfilePhotoRequest
	resp := new(kprofile.CreateProfilePhotoResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamsErr)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	res, err := rpc.ProfileSvc.CreateProfilePhoto(ctx, &kprofile.CreateProfilePhotoRequest{AccountId: c.MustGet(consts.AccountID).(string)})
	if err != nil {
		hlog.Error("rpc profile service err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, res)
}

// ClearProfilePhoto .
// @router /profile/mini/photo [DELETE]
func ClearProfilePhoto(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hprofile.ClearProfilePhotoRequest
	resp := new(kprofile.ClearProfilePhotoResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamsErr)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	res, err := rpc.ProfileSvc.ClearProfilePhoto(ctx, &kprofile.ClearProfilePhotoRequest{AccountId: c.MustGet(consts.AccountID).(string)})
	if err != nil {
		hlog.Error("rpc profile service err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetProfilePhoto .
// @router /profile/mini/photo [GET]
func GetProfilePhoto(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hprofile.GetProfilePhotoRequest
	resp := new(kprofile.GetProfilePhotoResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamsErr)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	res, err := rpc.ProfileSvc.GetProfilePhoto(ctx, &kprofile.GetProfilePhotoRequest{AccountId: c.MustGet(consts.AccountID).(string)})
	if err != nil {
		hlog.Error("rpc profile service err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, res)
}

// CompleteProfilePhoto .
// @router /profile/mini/complete [GET]
func CompleteProfilePhoto(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hprofile.CompleteProfilePhotoRequest
	resp := new(kprofile.CompleteProfilePhotoResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamsErr)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	res, err := rpc.ProfileSvc.CompleteProfilePhoto(ctx, &kprofile.CompleteProfilePhotoRequest{AccountId: c.MustGet(consts.AccountID).(string)})
	if err != nil {
		hlog.Error("rpc profile service err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, res)
}
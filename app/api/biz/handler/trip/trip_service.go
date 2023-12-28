// Code generated by hertz generator.

package trip

import (
	"context"
	"net/http"

	htrip "github.com/alimy/freecar/app/api/biz/model/trip"
	"github.com/alimy/freecar/app/api/internal/tool"
	"github.com/alimy/freecar/app/api/rpc"
	kbase "github.com/alimy/freecar/idle/auto/rpc/base"
	ktrip "github.com/alimy/freecar/idle/auto/rpc/trip"
	"github.com/alimy/freecar/library/core/consts"
	"github.com/alimy/freecar/library/core/errno"
	"github.com/alimy/freecar/library/core/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func DeleteTrip(ctx context.Context, c *app.RequestContext) {
	var err error
	var req htrip.DeleteTripRequest
	resp := new(ktrip.DeleteTripResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamsErr)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	res, err := rpc.TripSvc.DeleteTrip(ctx, &ktrip.DeleteTripRequest{Id: req.ID})
	if err != nil {
		hlog.Error("rpc trip service err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetAllTrips .
// @router /trip/admin/all [GET]
func GetAllTrips(ctx context.Context, c *app.RequestContext) {
	var err error
	var req htrip.GetAllTripsRequest
	resp := new(ktrip.GetAllTripsResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamsErr)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	res, err := rpc.TripSvc.GetAllTrips(ctx, &ktrip.GetAllTripsRequest{})
	if err != nil {
		hlog.Error("rpc trip service err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetSomeTrips .
// @router /trip/admin/some [GET]
func GetSomeTrips(ctx context.Context, c *app.RequestContext) {
	var err error
	var req htrip.GetSomeTripsRequest
	resp := new(ktrip.GetSomeTripsResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamsErr)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	res, err := rpc.TripSvc.GetSomeTrips(ctx, &ktrip.GetSomeTripsRequest{})
	if err != nil {
		hlog.Error("rpc trip service err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, res)
}

// CreateTrip .
// @router /trip/mini/trip [POST]
func CreateTrip(ctx context.Context, c *app.RequestContext) {
	var err error
	var req htrip.CreateTripRequest
	resp := new(ktrip.CreateTripResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamsErr)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	res, err := rpc.TripSvc.CreateTrip(ctx, &ktrip.CreateTripRequest{
		Start:     tool.ConvertTripLocation(req.Start),
		CarId:     req.CarID,
		AvatarUrl: req.AvatarURL,
		AccountId: c.MustGet(consts.AccountID).(string),
	})
	if err != nil {
		hlog.Error("rpc trip service err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetTrip .
// @router /trip/mini/trip [GET]
func GetTrip(ctx context.Context, c *app.RequestContext) {
	var err error
	var req htrip.GetTripRequest
	resp := new(ktrip.GetTripResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamsErr)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	res, err := rpc.TripSvc.GetTrip(ctx, &ktrip.GetTripRequest{
		Id:        req.ID,
		AccountId: c.MustGet(consts.AccountID).(string),
	})
	if err != nil {
		hlog.Error("rpc trip service err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetTrips .
// @router /trip/mini/trips [GET]
func GetTrips(ctx context.Context, c *app.RequestContext) {
	var err error
	var req htrip.GetTripsRequest
	resp := new(ktrip.GetTripsResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamsErr)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	res, err := rpc.TripSvc.GetTrips(ctx, &ktrip.GetTripsRequest{
		Status:    kbase.TripStatus(req.Status),
		AccountId: c.MustGet(consts.AccountID).(string),
	})
	if err != nil {
		hlog.Error("rpc trip service err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, res)
}

// UpdateTrip .
// @router /trip [PUT]
func UpdateTrip(ctx context.Context, c *app.RequestContext) {
	var err error
	var req htrip.UpdateTripRequest
	resp := new(ktrip.UpdateTripResponse)

	if err = c.BindAndValidate(&req); err != nil {
		resp.BaseResp = utils.BuildBaseResp(errno.ParamsErr)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	res, err := rpc.TripSvc.UpdateTrip(ctx, &ktrip.UpdateTripRequest{
		Id:        req.ID,
		Current:   (*kbase.Location)(req.Current),
		EndTrip:   req.EndTrip,
		AccountId: c.MustGet(consts.AccountID).(string),
	})
	if err != nil {
		hlog.Error("rpc trip service err", err)
		resp.BaseResp = utils.BuildBaseResp(errno.ServiceErr)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, res)
}

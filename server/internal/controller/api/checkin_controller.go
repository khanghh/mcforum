package api

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/controller/payload"
	"bbs-go/internal/service"
	"time"

	"bbs-go/common/dates"
	"bbs-go/pkg/web"

	"github.com/kataras/iris/v12"
)

type CheckinController struct {
	Ctx iris.Context
}

// PostCheckin Check in
func (c *CheckinController) PostCheckin() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if err := service.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}
	err := service.CheckInService.CheckIn(user.ID)
	if err == nil {
		return web.JsonSuccess()
	} else {
		return web.JsonError(err)
	}
}

// GetCheckin Get check-in info
func (c *CheckinController) GetCheckin() *web.JsonResult {
	user := service.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonSuccess()
	}
	checkIn := service.CheckInService.GetByUserId(user.ID)
	if checkIn != nil {
		today := dates.GetDay(time.Now())
		return web.NewRspBuilder(checkIn).
			Put("checkIn", checkIn.LatestDayName == today). // whether checked in today
			JsonResult()
	}
	return web.JsonSuccess()
}

// GetRank Get today's check-in leaderboard (earliest check-ins first)
func (c *CheckinController) GetRank() *web.JsonResult {
	list := cache.UserCache.GetCheckInRank()
	var itemList []map[string]interface{}
	for _, checkIn := range list {
		itemList = append(itemList, web.NewRspBuilder(checkIn).
			Put("user", payload.BuildUserInfoDefaultIfNull(checkIn.UserId)).
			Build())
	}
	return web.JsonData(itemList)
}

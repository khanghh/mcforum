package api

import (
	"bbs-go/internal/cache"
	"bbs-go/pkg/web"

	"github.com/kataras/iris/v12"
	"github.com/mcstatus-io/mcutil/v4/status"
)

type MCStatus struct {
	Status        string `json:"status"`
	Version       string `json:"version"`
	Protocol      int64  `json:"protocol"`
	PlayersOnline *int64 `json:"playersOnline"`
	PlayersMax    *int64 `json:"playersMax"`
}

type MCStatusController struct {
	Ctx   iris.Context
	cache *MCStatus
}

var mcVersions = map[int64]string{
	774: "1.21.11",
	773: "1.21.10",
	772: "1.21.8",
	771: "1.21.6",
	770: "1.21.5",
	769: "1.21.4",
	768: "1.21.3",
	767: "1.21.1",
	766: "1.20.6",
	765: "1.20.4",
	764: "1.20.2",
	763: "1.20.1",
	762: "1.19.4",
	761: "1.19.3",
	760: "1.19.2",
	759: "1.19",
	758: "1.18.2",
	757: "1.18.1",
	756: "1.17.1",
	755: "1.17",
	754: "1.16.5",
	753: "1.16.3",
	751: "1.16.2",
	736: "1.16.1",
	735: "1.16",
	578: "1.15.2",
	575: "1.15.1",
	573: "1.15",
	498: "1.14.4",
	477: "1.14",
	404: "1.13.2",
	393: "1.13",
	340: "1.12.2",
	316: "1.11.2",
	210: "1.10.2",
	110: "1.9.4",
	47:  "1.8.9",
	5:   "1.7.10",
	4:   "1.7.2",
}

func (c *MCStatusController) Get() *web.JsonResult {
	serverIP := cache.SysConfigCache.GetValue("mcServerIP")
	offlineStatus := MCStatus{Status: "offline"}
	if serverIP == "" {
		return web.JsonData(offlineStatus)
	}
	response, err := status.Modern(c.Ctx, serverIP, 25565)
	if err == nil {
		mcVersion := mcVersions[response.Version.Protocol]
		if mcVersion == "" {
			mcVersion = "Unknown"
		}
		c.cache = &MCStatus{
			Status:        "online",
			Version:       mcVersion,
			Protocol:      response.Version.Protocol,
			PlayersOnline: response.Players.Online,
			PlayersMax:    response.Players.Max,
		}
		return web.JsonData(c.cache)
	} else if c.cache != nil {
		c.cache.Status = "offline"
		return web.JsonData(c.cache)
	}
	return web.JsonData(offlineStatus)
}

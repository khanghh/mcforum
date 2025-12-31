package server

import (
	"bbs-go/internal/config"
	"bbs-go/internal/controller/api"
	"bbs-go/internal/middleware"
	"log/slog"
	"os"
	"strings"

	"bbs-go/pkg/web"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"

	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func NewServer() {
	conf := config.Instance()

	app := iris.New()
	app.Logger().SetLevel("info")
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Options{
		AllowedOrigins:   conf.AllowedOrigins,
		AllowCredentials: true,
		MaxAge:           600,
		AllowedMethods:   []string{iris.MethodOptions, iris.MethodHead, iris.MethodGet, iris.MethodPost, iris.MethodPut, iris.MethodPatch, iris.MethodDelete},
		AllowedHeaders:   []string{"*"},
	}))

	app.OnAnyErrorCode(func(ctx iris.Context) {
		path := ctx.Path()
		if strings.Contains(path, "/api/admin/") {
			err := ctx.JSON(web.JsonErrorCodeMsg(ctx.GetStatusCode(), "Http error"))
			slog.Error(err.Error(), slog.Any("err", err))
		}
	})

	// api
	apiRoute := NewMVCApplication(mvc.New(app.Party("/api")), kebabCasePathWordFunc)
	apiRoute.HandleError(middleware.ErrorHandler)
	apiRoute.Configure(func(m *mvc.Application) {
		apiRoute.Party("/login").Handle(new(api.LoginController))
		apiRoute.Party("/users/me").Handle(new(api.MeController))
		apiRoute.Party("/users").Handle(new(api.UsersController))
		apiRoute.Party("/topics").Handle(new(api.TopicsController))
		apiRoute.Party("/feeds").Handle(new(api.FeedsController))
		apiRoute.Party("/forums").Handle(new(api.ForumsController))
		apiRoute.Party("/comments").Handle(new(api.CommentsController))
		apiRoute.Party("/tags").Handle(new(api.TagsController))
		apiRoute.Party("/favorites").Handle(new(api.FavoritesController))
		apiRoute.Party("/likes").Handle(new(api.LikeController))
		apiRoute.Party("/checkin").Handle(new(api.CheckinController))
		apiRoute.Party("/config").Handle(new(api.ConfigController))
		apiRoute.Party("/upload").Handle(new(api.UploadController))
		apiRoute.Party("/link").Handle(new(api.LinkController))
		apiRoute.Party("/captcha").Handle(new(api.CaptchaController))
		apiRoute.Party("/search").Handle(new(api.SearchController))
		// apiRoute.Party("/fans").Handle(new(api.FansController))
		// apiRoute.Party("/user-report").Handle(new(api.UserReportController))
	})

	// admin
	// mvc.Configure(app.Party("/api/admin"), func(m *mvc.Application) {
	// 	m.Router.Use(middleware.AdminAuth)
	// 	m.Party("/common").Handle(new(admin.CommonController))
	// 	m.Party("/user").Handle(new(admin.UserController))
	// 	// m.Party("/tag").Handle(new(admin.TagController))
	// 	m.Party("/comment").Handle(new(admin.CommentController))
	// 	m.Party("/favorite").Handle(new(admin.FavoriteController))
	// 	m.Party("/topics").Handle(new(admin.TopicController))
	// 	m.Party("/topic-node").Handle(new(admin.ForumController))
	// 	m.Party("/sys-config").Handle(new(admin.SysConfigController))
	// 	m.Party("/link").Handle(new(admin.LinkController))
	// 	m.Party("/user-score-log").Handle(new(admin.UserScoreLogController))
	// 	m.Party("/operate-log").Handle(new(admin.OperateLogController))
	// 	m.Party("/user-report").Handle(new(admin.UserReportController))
	// 	m.Party("/forbidden-word").Handle(new(admin.ForbiddenWordController))

	// 	m.Party("/role").Handle(new(admin.RoleController))
	// 	m.Party("/menu").Handle(new(admin.MenuController))
	// })
	// for _, route := range app.GetRoutes() {
	// 	if route.Method != iris.MethodOptions {
	// 		fmt.Printf("%s\t\t%s\n", route.Method, route.Path)
	// 	}
	// }

	if err := app.Listen(":"+conf.Port,
		iris.WithConfiguration(iris.Configuration{
			DisableStartupLog:                 false,
			DisableInterruptHandler:           false,
			DisablePathCorrection:             false,
			EnablePathEscape:                  false,
			FireMethodNotAllowed:              false,
			DisableBodyConsumptionOnUnmarshal: false,
			DisableAutoFireStatusCode:         false,
			EnableOptimizations:               true,
			TimeFormat:                        "2006-01-02 15:04:05",
			Charset:                           "UTF-8",
		}),
	); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		os.Exit(-1)
	}
}

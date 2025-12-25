package middleware

import (
	"bbs-go/common/urls"
	"bbs-go/internal/errs"
	"bbs-go/internal/locale"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/service"
	"errors"
	"fmt"

	"bbs-go/pkg/web"

	"github.com/kataras/iris/v12"
)

var (
	config = []PathRole{
		{Pattern: "/api/admin/sys-config/**", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/user/create", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/user/update", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/topic-node/create", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/topic-node/update", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/tag/create", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/tag/update", Roles: []string{constants.RoleOwner}},
		{Pattern: "/api/admin/**", Roles: []string{constants.RoleOwner, constants.RoleAdmin}},
	}
	antPathMatcher = urls.NewAntPathMatcher()
)

// AdminAuth admin permission
func AdminAuth(ctx iris.Context) {
	roles := getPathRoles(ctx)

	// no role required to access
	if len(roles) == 0 {
		return
	}

	user := service.UserTokenService.GetCurrent(ctx)
	if user == nil {
		notLogin(ctx)
		return
	}
	if !user.HasAnyRole(roles...) {
		fmt.Println("noPermission")
		noPermission(ctx)
		return
	}

	ctx.Next()
}

// getPathRoles get roles required for this path
func getPathRoles(ctx iris.Context) []string {
	p := ctx.Path()
	for _, pathRole := range config {
		if antPathMatcher.Match(pathRole.Pattern, p) {
			return pathRole.Roles
		}
	}
	return nil
}

// notLogin not logged in response
func notLogin(ctx iris.Context) {
	_ = ctx.JSON(web.JsonError(errs.NotLogin))
	ctx.StopExecution()
}

// noPermission no permission response
func noPermission(ctx iris.Context) {
	_ = ctx.JSON(web.JsonErrorCode(iris.StatusForbidden, errors.New(locale.T("errors.permission_denied"))))
	ctx.StopExecution()
}

type PathRole struct {
	Pattern string   // path pattern
	Roles   []string // roles
}

package server

import (
	"fmt"
	"strings"

	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/mvc"
)

func kebabCasePathWordFunc(path string, w string, index int) string {
	fmt.Printf("path:%s, w:%s, index:%d\n", path, w, index)
	if strings.HasSuffix(path, "}") {
		return path + "/" + strings.ToLower(w)
	}
	if index == 0 {
		return path + strings.ToLower(w)
	}
	if strings.HasSuffix(path, "_") {
		return path[:len(path)-1] + "/" + strings.ToLower(w)
	}
	if strings.Contains(path, "_") {
		path = strings.ReplaceAll(path, "_", "/")
	}

	return path + "-" + strings.ToLower(w)
}

type MVCApplication struct {
	*mvc.Application
	customPathWordFunc mvc.CustomPathWordFunc
}

func NewMVCApplication(app *mvc.Application, wordFunc mvc.CustomPathWordFunc) *MVCApplication {
	return &MVCApplication{
		Application:        app.SetCustomPathWordFunc(wordFunc),
		customPathWordFunc: wordFunc,
	}
}

func (app *MVCApplication) Clone(party router.Party) *MVCApplication {
	return NewMVCApplication(app.Application.Clone(party), app.customPathWordFunc)
}

func (app *MVCApplication) Party(relativePath string, middleware ...context.Handler) *MVCApplication {
	return app.Clone(app.Router.Party(relativePath, middleware...))
}

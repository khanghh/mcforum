package server

import (
	"fmt"
	"strings"

	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/mvc"
)

func kebabCasePathWordFunc(path string, w string, index int) string {
	if index > 0 {
		path += "-" + strings.ToLower(w)
	} else {
		path += strings.ToLower(w)
	}
	// if path[len(path)-1] != '/' {
	// 	path += "/"
	// }
	fmt.Println("path", path)
	return path
	// if index == 0 {
	// 	return path + strings.ToLower(word)
	// }
	// return path + "-" + strings.ToLower(word)
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

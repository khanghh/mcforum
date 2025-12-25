package params

import (
	"strings"

	"bbs-go/common/strs"
	"bbs-go/common/strs/strcase"
	"bbs-go/sqls"

	"github.com/kataras/iris/v12"
)

type QueryOp string

const (
	Eq       QueryOp = "eq"
	Gt       QueryOp = "gt"
	Lt       QueryOp = "lt"
	Gte      QueryOp = "gte"
	Lte      QueryOp = "lte"
	Like     QueryOp = "like"
	In       QueryOp = "in"
	Starting QueryOp = "starting"
	Ending   QueryOp = "ending"
)

type QueryFilter struct {
	ParamName    string                     // Request parameter name
	Op           QueryOp                    // Operator
	ColumnName   string                     // Column name
	ValueWrapper func(origin string) string // Value decorator, can
}

func NewPagedSqlCnd(ctx iris.Context, filters ...QueryFilter) *sqls.Cnd {
	cnd := NewSqlCnd(ctx, filters...)
	p := GetPaging(ctx)
	cnd.Page(p.Page, p.Limit)
	return cnd
}

func NewSqlCnd(ctx iris.Context, filters ...QueryFilter) *sqls.Cnd {
	cnd := sqls.NewCnd()
	for _, filter := range filters {
		var (
			columnName = filter.ColumnName
			paramValue = ctx.FormValue(filter.ParamName)
		)
		if strs.IsBlank(paramValue) {
			continue
		}
		if filter.ValueWrapper != nil {
			paramValue = filter.ValueWrapper(paramValue)
		}
		if strs.IsBlank(string(filter.Op)) {
			filter.Op = Eq
		}
		if strs.IsBlank(columnName) {
			columnName = strcase.ToSnake(filter.ParamName)
		}
		if filter.Op == Eq {
			cnd.Eq(columnName, paramValue)
		} else if filter.Op == Gt {
			cnd.Gt(columnName, paramValue)
		} else if filter.Op == Lt {
			cnd.Lt(columnName, paramValue)
		} else if filter.Op == Gte {
			cnd.Gte(columnName, paramValue)
		} else if filter.Op == Lte {
			cnd.Lte(columnName, paramValue)
		} else if filter.Op == Like {
			cnd.Like(columnName, paramValue)
		} else if filter.Op == Starting {
			cnd.Starting(columnName, paramValue)
		} else if filter.Op == Ending {
			cnd.Ending(columnName, paramValue)
		} else if filter.Op == In {
			ss := strings.Split(paramValue, ",")
			cnd.In(columnName, ss)
		}
	}
	return cnd
}

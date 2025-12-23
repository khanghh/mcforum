package web

import (
	"bbs-go/common/structs"
	"bbs-go/sqls"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kataras/iris/v12"
)

type jsonError interface {
	Code() int
	Error() string
}

type JsonResult struct {
	StatusCode int
	Data       interface{}
	Error      error
	Properties map[string]interface{}
}

type jsonErrorMarshaling struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//	type jsonResultMarshaling struct {
//		ApiVersion string                 `json:"apiVersion,omitempty"`
//		Context    string                 `json:"context,omitempty"`
//		Id         string                 `json:"id,omitempty"`
//		Method     string                 `json:"method,omitempty"`
//		Params     map[string]interface{} `json:"params,omitempty"`
//		Data       interface{}            `json:"data,omitempty"`
//		Error      *jsonErrorMarshaling   `json:"error,omitempty"`
//	}
func (r *JsonResult) SetProperty(key string, val interface{}) *JsonResult {
	if r.Properties == nil {
		r.Properties = make(map[string]interface{})
	}
	r.Properties[key] = val
	return r
}

func (r *JsonResult) MarshalJSON() ([]byte, error) {
	ret := make(map[string]interface{})
	for key, val := range r.Properties {
		ret[key] = val
	}
	if r.Data != nil {
		ret["data"] = r.Data
	}
	if r.Error != nil {
		if jsonErr, ok := r.Error.(jsonError); ok {
			ret["error"] = jsonErrorMarshaling{
				Code:    jsonErr.Code(),
				Message: jsonErr.Error(),
			}
		} else {
			statusCode := r.StatusCode
			if statusCode == 0 {
				statusCode = iris.StatusBadRequest
			}
			ret["error"] = jsonErrorMarshaling{
				Code:    statusCode,
				Message: r.Error.Error(),
			}
		}
	}
	if result, ok := r.Data.(*CursorResult); ok {
		fmt.Println("ret", result.Items)
	}
	return json.Marshal(ret)
}

func (r *JsonResult) Dispatch(ctx iris.Context) {
	if r.StatusCode != 0 {
		ctx.StatusCode(r.StatusCode)
	} else if jsonErr, ok := r.Error.(jsonError); ok {
		ctx.StatusCode(jsonErr.Code())
	}
	ctx.JSON(r)
}

func Json(statusCode int, err error, data interface{}) *JsonResult {
	if err != nil {
		return &JsonResult{
			StatusCode: statusCode,
			Error:      err,
			Data:       data,
		}
	}
	return &JsonResult{
		StatusCode: statusCode,
		Data:       data,
	}
}

func JsonData(data interface{}) *JsonResult {
	return &JsonResult{
		StatusCode: iris.StatusOK,
		Data:       data,
	}
}

func JsonPageData(results interface{}, page *sqls.Paging) *JsonResult {
	return JsonData(&PageResult{
		Results: results,
		Page:    page,
	})
}

func JsonCursorData(items interface{}, cursor string, hasMore bool) *JsonResult {
	return JsonData(&CursorResult{
		Items:   items,
		Cursor:  cursor,
		HasMore: hasMore,
	})
}

func JsonSuccess() *JsonResult {
	return &JsonResult{StatusCode: iris.StatusOK}
}

func JsonError(err error) *JsonResult {
	return &JsonResult{Error: err}
}

func JsonErrorMsg(msg string) *JsonResult {
	return JsonError(errors.New(msg))
}

func JsonErrorCode(statusCode int, err error) *JsonResult {
	return &JsonResult{
		StatusCode: statusCode,
		Error:      err,
	}
}

func JsonErrorCodeMsg(statusCode int, msg string) *JsonResult {
	return JsonErrorCode(statusCode, errors.New(msg))
}

type RspBuilder struct {
	Data map[string]interface{}
}

func NewEmptyRspBuilder() *RspBuilder {
	return &RspBuilder{Data: make(map[string]interface{})}
}

func NewRspBuilder(obj interface{}) *RspBuilder {
	return NewRspBuilderExcludes(obj)
}

func NewRspBuilderExcludes(obj interface{}, excludes ...string) *RspBuilder {
	return &RspBuilder{Data: structs.StructToMap(obj, excludes...)}
}

func (builder *RspBuilder) Put(key string, value interface{}) *RspBuilder {
	builder.Data[key] = value
	return builder
}

func (builder *RspBuilder) Build() map[string]interface{} {
	return builder.Data
}

func (builder *RspBuilder) JsonResult() *JsonResult {
	return JsonData(builder.Data)
}

func ConvertList[T any](results []T, conv func(item T) map[string]interface{}) (list []map[string]interface{}) {
	for _, item := range results {
		if ret := conv(item); ret != nil {
			list = append(list, ret)
		}
	}
	return
}

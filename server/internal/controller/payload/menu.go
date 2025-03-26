package payload

import "bbs-go/internal/model"

type MenuResponse struct {
	Id         int64  `json:"id"`
	ParentId   *int64 `json:"parentId"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	SortNo     int    `json:"sortNo"`
	Status     int    `json:"status"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
}

type MenuTreeResponse struct {
	MenuResponse
	Level    int                `json:"level"`
	Children []MenuTreeResponse `json:"children"`
}

type TreeNode struct {
	Id       int64      `json:"id"`
	Key      int64      `json:"key"`
	Title    string     `json:"title"`
	Children []TreeNode `json:"children"`
}

func BuildMenu(element *model.Menu) MenuResponse {
	item := MenuResponse{
		Id:         element.Id,
		Name:       element.Name,
		Title:      element.Title,
		Icon:       element.Icon,
		Path:       element.Path,
		SortNo:     element.SortNo,
		Status:     element.Status,
		CreateTime: element.CreateTime,
		UpdateTime: element.UpdateTime,
	}
	if element.ParentId > 0 {
		item.ParentId = &element.ParentId
	}
	return item
}

func BuildMenuTree(parentId int64, list []model.Menu) (ret []MenuTreeResponse) {
	return _BuildMenuTree(parentId, 1, list)
}

func _BuildMenuTree(parentId int64, level int, list []model.Menu) (ret []MenuTreeResponse) {
	for _, element := range list {
		if element.ParentId == parentId {
			menu := BuildMenu(&element)
			ret = append(ret, MenuTreeResponse{
				MenuResponse: menu,
				Level:        level,
				Children:     _BuildMenuTree(element.Id, level+1, list),
			})
		}
	}
	return
}

func BuildMenuSimpleTree(parentId int64, list []model.Menu) (ret []TreeNode) {
	return _BuildMenuSimpleTree(parentId, list)
}

func _BuildMenuSimpleTree(parentId int64, list []model.Menu) (ret []TreeNode) {
	for _, element := range list {
		if element.ParentId == parentId {
			ret = append(ret, TreeNode{
				Id:       element.Id,
				Key:      element.Id,
				Title:    element.Title,
				Children: _BuildMenuSimpleTree(element.Id, list),
			})
		}
	}
	return
}

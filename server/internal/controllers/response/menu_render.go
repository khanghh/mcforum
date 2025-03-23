package response

import "bbs-go/internal/models"

func BuildMenu(element *models.Menu) MenuResponse {
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

func BuildMenuTree(parentId int64, list []models.Menu) (ret []MenuTreeResponse) {
	return _BuildMenuTree(parentId, 1, list)
}

func _BuildMenuTree(parentId int64, level int, list []models.Menu) (ret []MenuTreeResponse) {
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

func BuildMenuSimpleTree(parentId int64, list []models.Menu) (ret []TreeNode) {
	return _BuildMenuSimpleTree(parentId, list)
}

func _BuildMenuSimpleTree(parentId int64, list []models.Menu) (ret []TreeNode) {
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

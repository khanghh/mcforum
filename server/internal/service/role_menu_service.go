package service

import (
	"bbs-go/internal/model"
	"bbs-go/internal/repository"

	"bbs-go/common/arrs"
	"bbs-go/common/dates"
	"bbs-go/sqls"
	"bbs-go/web/params"

	"gorm.io/gorm"
)

var RoleMenuService = newRoleMenuService()

func newRoleMenuService() *roleMenuService {
	return &roleMenuService{}
}

type roleMenuService struct {
}

func (s *roleMenuService) Get(id int64) *model.RoleMenu {
	return repository.RoleMenuRepository.Get(sqls.DB(), id)
}

func (s *roleMenuService) Take(where ...interface{}) *model.RoleMenu {
	return repository.RoleMenuRepository.Take(sqls.DB(), where...)
}

func (s *roleMenuService) Find(cnd *sqls.Cnd) []model.RoleMenu {
	return repository.RoleMenuRepository.Find(sqls.DB(), cnd)
}

func (s *roleMenuService) FindOne(cnd *sqls.Cnd) *model.RoleMenu {
	return repository.RoleMenuRepository.FindOne(sqls.DB(), cnd)
}

func (s *roleMenuService) FindPageByParams(params *params.QueryParams) (list []model.RoleMenu, paging *sqls.Paging) {
	return repository.RoleMenuRepository.FindPageByParams(sqls.DB(), params)
}

func (s *roleMenuService) FindPageByCnd(cnd *sqls.Cnd) (list []model.RoleMenu, paging *sqls.Paging) {
	return repository.RoleMenuRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *roleMenuService) Count(cnd *sqls.Cnd) int64 {
	return repository.RoleMenuRepository.Count(sqls.DB(), cnd)
}

func (s *roleMenuService) Create(t *model.RoleMenu) error {
	return repository.RoleMenuRepository.Create(sqls.DB(), t)
}

func (s *roleMenuService) Update(t *model.RoleMenu) error {
	return repository.RoleMenuRepository.Update(sqls.DB(), t)
}

func (s *roleMenuService) Updates(id int64, columns map[string]interface{}) error {
	return repository.RoleMenuRepository.Updates(sqls.DB(), id, columns)
}

func (s *roleMenuService) UpdateColumn(id int64, name string, value interface{}) error {
	return repository.RoleMenuRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *roleMenuService) Delete(id int64) {
	repository.RoleMenuRepository.Delete(sqls.DB(), id)
}

func (s *roleMenuService) GetByRole(roleId int64) []model.RoleMenu {
	return s.Find(sqls.NewCnd().Eq("role_id", roleId))
}

func (s *roleMenuService) GetMenuIdsByRoles(roleIds []int64) (menuIds []int64) {
	list := s.Find(sqls.NewCnd().In("role_id", roleIds))
	for _, element := range list {
		menuIds = append(menuIds, element.MenuId)
	}
	return
}

func (s *roleMenuService) GetMenuIdsByRole(roleId int64) (menuIds []int64) {
	list := s.GetByRole(roleId)
	for _, element := range list {
		menuIds = append(menuIds, element.MenuId)
	}
	return
}

func (s *roleMenuService) SaveRoleMenus(roleId int64, menuIds []int64) error {
	currentMenuIds := s.GetMenuIdsByRole(roleId)
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		var (
			addIds []int64 // 本次需要新增的
			delIds []int64 // 本次需要删除的
		)
		for _, menuId := range menuIds {
			if !arrs.Contains(currentMenuIds, menuId) {
				addIds = append(addIds, menuId)
			}
		}
		for _, menuId := range currentMenuIds {
			if !arrs.Contains(menuIds, menuId) {
				delIds = append(delIds, menuId)
			}
		}

		for _, menuId := range addIds {
			if err := repository.RoleMenuRepository.Create(tx, &model.RoleMenu{
				RoleId:     roleId,
				MenuId:     menuId,
				CreateTime: dates.NowTimestamp(),
			}); err != nil {
				return err
			}
		}
		for _, menuId := range delIds {
			if err := tx.Delete(&model.RoleMenu{}, "role_id = ? and menu_id = ?", roleId, menuId).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

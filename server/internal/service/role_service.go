package service

import (
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/repository"

	"bbs-go/pkg/web/params"
	"bbs-go/sqls"

	"gorm.io/gorm"
)

var RoleService = newRoleService()

func newRoleService() *roleService {
	return &roleService{}
}

type roleService struct {
}

func (s *roleService) Get(id int64) *model.Role {
	return repository.RoleRepository.Get(sqls.DB(), id)
}

func (s *roleService) Take(where ...interface{}) *model.Role {
	return repository.RoleRepository.Take(sqls.DB(), where...)
}

func (s *roleService) Find(cnd *sqls.Cnd) []model.Role {
	return repository.RoleRepository.Find(sqls.DB(), cnd)
}

func (s *roleService) FindOne(cnd *sqls.Cnd) *model.Role {
	return repository.RoleRepository.FindOne(sqls.DB(), cnd)
}

func (s *roleService) FindPageByParams(params *params.QueryParams) (list []model.Role, paging *sqls.Paging) {
	return repository.RoleRepository.FindPageByParams(sqls.DB(), params)
}

func (s *roleService) FindPageByCnd(cnd *sqls.Cnd) (list []model.Role, paging *sqls.Paging) {
	return repository.RoleRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *roleService) Count(cnd *sqls.Cnd) int64 {
	return repository.RoleRepository.Count(sqls.DB(), cnd)
}

func (s *roleService) Create(t *model.Role) error {
	return repository.RoleRepository.Create(sqls.DB(), t)
}

func (s *roleService) Update(t *model.Role) error {
	return repository.RoleRepository.Update(sqls.DB(), t)
}

func (s *roleService) Updates(id int64, columns map[string]interface{}) error {
	return repository.RoleRepository.Updates(sqls.DB(), id, columns)
}

func (s *roleService) UpdateColumn(id int64, name string, value interface{}) error {
	return repository.RoleRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *roleService) Delete(id int64) {
	repository.RoleRepository.Delete(sqls.DB(), id)
}

func (s *roleService) GetByCode(code string) *model.Role {
	return s.FindOne(sqls.NewCnd().Eq("code", code))
}

func (s *roleService) GetNextSortNo() int {
	if max := s.FindOne(sqls.NewCnd().Eq("status", constants.StatusActive).Desc("sort_no")); max != nil {
		return max.SortNo + 1
	}
	return 0
}

func (s *roleService) UpdateSort(ids []int64) error {
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := repository.RoleRepository.UpdateColumn(tx, id, "sort_no", i); err != nil {
				return err
			}
		}
		return nil
	})
}

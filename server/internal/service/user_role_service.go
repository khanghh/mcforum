package service

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/model"
	"bbs-go/internal/repository"
	"strings"

	"bbs-go/common/dates"
	"bbs-go/sqls"
	"bbs-go/web/params"

	"gorm.io/gorm"
)

var UserRoleService = newUserRoleService()

func newUserRoleService() *userRoleService {
	return &userRoleService{}
}

type userRoleService struct {
}

func (s *userRoleService) Get(id int64) *model.UserRole {
	return repository.UserRoleRepository.Get(sqls.DB(), id)
}

func (s *userRoleService) Take(where ...interface{}) *model.UserRole {
	return repository.UserRoleRepository.Take(sqls.DB(), where...)
}

func (s *userRoleService) Find(cnd *sqls.Cnd) []model.UserRole {
	return repository.UserRoleRepository.Find(sqls.DB(), cnd)
}

func (s *userRoleService) FindOne(cnd *sqls.Cnd) *model.UserRole {
	return repository.UserRoleRepository.FindOne(sqls.DB(), cnd)
}

func (s *userRoleService) FindPageByParams(params *params.QueryParams) (list []model.UserRole, paging *sqls.Paging) {
	return repository.UserRoleRepository.FindPageByParams(sqls.DB(), params)
}

func (s *userRoleService) FindPageByCnd(cnd *sqls.Cnd) (list []model.UserRole, paging *sqls.Paging) {
	return repository.UserRoleRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *userRoleService) Count(cnd *sqls.Cnd) int64 {
	return repository.UserRoleRepository.Count(sqls.DB(), cnd)
}

func (s *userRoleService) Delete(id int64) {
	repository.UserRoleRepository.Delete(sqls.DB(), id)
}

func (s *userRoleService) UpdateUserRoles(userId int64, roleIds []int64) error {
	err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		var roles []model.Role
		if len(roleIds) > 0 {
			roles = repository.RoleRepository.Find(tx, sqls.NewCnd().In("id", roleIds))
		}

		var roleCodes []string
		for _, role := range roles {
			roleCodes = append(roleCodes, role.Code)
		}

		if err := tx.Delete(&model.UserRole{}, "user_id = ?", userId).Error; err != nil {
			return err
		}
		if len(roles) == 0 {
			return repository.UserRepository.UpdateColumn(tx, userId, "roles", "")
		} else {
			for _, role := range roles {
				if err := repository.UserRoleRepository.Create(tx, &model.UserRole{
					UserId:     userId,
					RoleId:     role.Id,
					CreateTime: dates.NowTimestamp(),
				}); err != nil {
					return err
				}
			}
			return repository.UserRepository.UpdateColumn(tx, userId, "roles", strings.Join(roleCodes, ","))
		}
	})
	if err != nil {
		return err
	}
	cache.UserCache.Invalidate(userId)
	return nil
}

func (s *userRoleService) GetUserRoleIds(userId int64) (roleIds []int64) {
	list := s.Find(sqls.NewCnd().Eq("user_id", userId))
	for _, userRole := range list {
		roleIds = append(roleIds, userRole.RoleId)
	}
	return roleIds
}

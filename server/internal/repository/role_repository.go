package repository

import (
	"bbs-go/internal/model"

	"bbs-go/sqls"
	"bbs-go/web/params"

	"gorm.io/gorm"
)

var RoleRepository = newRoleRepository()

func newRoleRepository() *roleRepository {
	return &roleRepository{}
}

type roleRepository struct {
}

func (r *roleRepository) Get(db *gorm.DB, id int64) *model.Role {
	ret := &model.Role{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *roleRepository) Take(db *gorm.DB, where ...interface{}) *model.Role {
	ret := &model.Role{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *roleRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.Role) {
	cnd.Find(db, &list)
	return
}

func (r *roleRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.Role {
	ret := &model.Role{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *roleRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []model.Role, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *roleRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []model.Role, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Role{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *roleRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (list []model.Role) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *roleRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *roleRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &model.Role{})
}

func (r *roleRepository) Create(db *gorm.DB, t *model.Role) (err error) {
	err = db.Create(t).Error
	return
}

func (r *roleRepository) Update(db *gorm.DB, t *model.Role) (err error) {
	err = db.Save(t).Error
	return
}

func (r *roleRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Role{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *roleRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Role{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *roleRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.Role{}, "id = ?", id)
}

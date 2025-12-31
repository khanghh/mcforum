package repository

import (
	"bbs-go/sqls"

	"gorm.io/gorm"

	"bbs-go/internal/model"
)

var UploadRepository = newUploadRepository()

func newUploadRepository() *uploadRepository {
	return &uploadRepository{}
}

type uploadRepository struct {
}

func (r *uploadRepository) Get(db *gorm.DB, id int64) *model.Upload {
	ret := &model.Upload{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *uploadRepository) Take(db *gorm.DB, where ...interface{}) *model.Upload {
	ret := &model.Upload{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *uploadRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.Upload) {
	cnd.Find(db, &list)
	return
}

func (r *uploadRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.Upload {
	ret := &model.Upload{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *uploadRepository) Create(db *gorm.DB, t *model.Upload) (err error) {
	err = db.Create(t).Error
	return
}

func (r *uploadRepository) Update(db *gorm.DB, t *model.Upload) (err error) {
	err = db.Save(t).Error
	return
}

func (r *uploadRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Link{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *uploadRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Upload{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *uploadRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.Upload{}, "id = ?", id)
}

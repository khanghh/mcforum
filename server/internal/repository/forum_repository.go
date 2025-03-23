package repository

import (
	"bbs-go/sqls"
	"bbs-go/web/params"

	"gorm.io/gorm"

	"bbs-go/internal/model"
)

var ForumRepository = newForumRepository()

func newForumRepository() *forumRepository {
	return &forumRepository{}
}

type forumRepository struct {
}

func (r *forumRepository) Get(db *gorm.DB, id int64) *model.Forum {
	ret := &model.Forum{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *forumRepository) Take(db *gorm.DB, where ...interface{}) *model.Forum {
	ret := &model.Forum{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *forumRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.Forum) {
	cnd.Find(db, &list)
	return
}

func (r *forumRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.Forum {
	ret := &model.Forum{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *forumRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []model.Forum, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *forumRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []model.Forum, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Forum{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *forumRepository) Create(db *gorm.DB, t *model.Forum) (err error) {
	err = db.Create(t).Error
	return
}

func (r *forumRepository) Update(db *gorm.DB, t *model.Forum) (err error) {
	err = db.Save(t).Error
	return
}

func (r *forumRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Forum{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *forumRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Forum{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *forumRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.Forum{}, "id = ?", id)
}

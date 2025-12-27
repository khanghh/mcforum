package service

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/event"
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/repository"

	"bbs-go/common/dates"
	"bbs-go/pkg/web/params"
	"bbs-go/sqls"

	"github.com/emirpasic/gods/sets/hashset"
	"gorm.io/gorm"
)

var UserFollowService = newUserFollowService()

func newUserFollowService() *userFollowService {
	return &userFollowService{}
}

type userFollowService struct {
}

func (s *userFollowService) Get(id int64) *model.UserFollow {
	return repository.UserFollowRepository.Get(sqls.DB(), id)
}

func (s *userFollowService) Take(where ...interface{}) *model.UserFollow {
	return repository.UserFollowRepository.Take(sqls.DB(), where...)
}

func (s *userFollowService) Find(cnd *sqls.Cnd) []model.UserFollow {
	return repository.UserFollowRepository.Find(sqls.DB(), cnd)
}

func (s *userFollowService) FindOne(cnd *sqls.Cnd) *model.UserFollow {
	return repository.UserFollowRepository.FindOne(sqls.DB(), cnd)
}

func (s *userFollowService) FindPageByParams(params *params.QueryParams) (list []model.UserFollow, paging *sqls.Paging) {
	return repository.UserFollowRepository.FindPageByParams(sqls.DB(), params)
}

func (s *userFollowService) FindPageByCnd(cnd *sqls.Cnd) (list []model.UserFollow, paging *sqls.Paging) {
	return repository.UserFollowRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *userFollowService) Count(cnd *sqls.Cnd) int64 {
	return repository.UserFollowRepository.Count(sqls.DB(), cnd)
}

func (s *userFollowService) Create(t *model.UserFollow) error {
	return repository.UserFollowRepository.Create(sqls.DB(), t)
}

func (s *userFollowService) Update(t *model.UserFollow) error {
	return repository.UserFollowRepository.Update(sqls.DB(), t)
}

func (s *userFollowService) Updates(id int64, columns map[string]interface{}) error {
	return repository.UserFollowRepository.Updates(sqls.DB(), id, columns)
}

func (s *userFollowService) UpdateColumn(id int64, name string, value interface{}) error {
	return repository.UserFollowRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *userFollowService) Delete(id int64) {
	repository.UserFollowRepository.Delete(sqls.DB(), id)
}

func (s *userFollowService) Follow(userId, otherId int64) error {
	if userId == otherId {
		// Following oneself: no-op.
		// return errors.New("Cannot follow yourself")
		return nil
	}

	if s.IsFollowing(userId, otherId) {
		return nil
	}

	err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		// If the other also follows me, update status to mutual follow
		otherFollowed := tx.Exec("update t_user_follow set status = ? where user_id = ? and other_id = ?",
			constants.FollowStatusBoth, otherId, userId).RowsAffected > 0
		status := constants.FollowStatusFollow
		if otherFollowed {
			status = constants.FollowStatusBoth
		}

		if err := repository.UserRepository.Updates(tx, userId, map[string]interface{}{
			"follow_count": gorm.Expr("follow_count + 1"),
		}); err != nil {
			return err
		}
		cache.UserCache.Invalidate(userId)

		if err := repository.UserRepository.Updates(tx, otherId, map[string]interface{}{
			"fans_count": gorm.Expr("fans_count + 1"),
		}); err != nil {
			return err
		}
		cache.UserCache.Invalidate(otherId)

		return repository.UserFollowRepository.Create(tx, &model.UserFollow{
			UserID:     userId,
			OtherID:    otherId,
			Status:     status,
			CreateTime: dates.NowTimestamp(),
		})
	})
	if err != nil {
		return err
	}

	// send mq message
	event.Send(event.FollowEvent{
		UserId:  userId,
		OtherId: otherId,
	})
	return nil
}

func (s *userFollowService) UnFollow(userId, otherId int64) error {
	if userId == otherId {
		// Following oneself: no-op.
		return nil
	}
	if !s.IsFollowing(userId, otherId) {
		return nil
	}
	err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		success := tx.Where("user_id = ? and other_id = ?", userId, otherId).Delete(model.UserFollow{}).RowsAffected > 0
		if success {
			tx.Exec("update t_user_follow set status = ? where user_id = ? and other_id = ?",
				constants.FollowStatusFollow, otherId, userId)
		}

		if err := tx.Model(&model.User{}).Where("id = ? and follow_count > 0", userId).Updates(map[string]interface{}{
			"follow_count": gorm.Expr("follow_count - 1"),
		}).Error; err != nil {
			return err
		}
		cache.UserCache.Invalidate(userId)

		if err := tx.Model(&model.User{}).Where("id = ? and fans_count > 0", otherId).Updates(map[string]interface{}{
			"fans_count": gorm.Expr("fans_count - 1"),
		}).Error; err != nil {
			return err
		}
		cache.UserCache.Invalidate(otherId)

		return nil
	})
	if err != nil {
		return err
	}

	// send mq message
	event.Send(event.UnfollowEvent{
		UserId:  userId,
		OtherId: otherId,
	})
	return nil
}

// GetFollowers Followers list
func (s *userFollowService) GetFollowers(userId int64, cursor int64, limit int) (itemList []int64, nextCursor int64, hasMore bool) {
	cnd := sqls.NewCnd().Eq("other_id", userId)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	cnd.Desc("id").Limit(limit)
	list := repository.UserFollowRepository.Find(sqls.DB(), cnd)

	if len(list) > 0 {
		nextCursor = list[len(list)-1].ID
		hasMore = len(list) >= limit
		for _, e := range list {
			itemList = append(itemList, e.UserID)
		}
	} else {
		nextCursor = cursor
	}
	return
}

// GetFollowing Following list
func (s *userFollowService) GetFollowing(userId int64, cursor int64, limit int) (itemList []int64, nextCursor int64, hasMore bool) {
	cnd := sqls.NewCnd().Eq("user_id", userId)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	cnd.Desc("id").Limit(limit)
	list := repository.UserFollowRepository.Find(sqls.DB(), cnd)

	if len(list) > 0 {
		nextCursor = list[len(list)-1].ID
		hasMore = len(list) >= limit
		for _, e := range list {
			itemList = append(itemList, e.OtherID)
		}
	} else {
		nextCursor = cursor
	}
	return
}

// ScanFollowers Scan followers
func (s *userFollowService) ScanFollowers(userId int64, handle func(fansId int64)) {
	var cursor int64 = 0
	for {
		list := s.Find(sqls.NewCnd().Eq("other_id", userId).Gt("id", cursor).Asc("id").Limit(100))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		for _, item := range list {
			handle(item.UserID)
		}
	}
}

// ScanFollowing scan following users of the given userId
func (s *userFollowService) ScanFollowing(userId int64, handle func(followUserId int64)) {
	var cursor int64 = 0
	for {
		list := s.Find(sqls.NewCnd().Eq("user_id", userId).Gt("id", cursor).Asc("id").Limit(100))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		for _, item := range list {
			handle(item.OtherID)
		}
	}
}

func (s *userFollowService) IsFollowing(userId, otherId int64) bool {
	if userId == otherId {
		return false
	}
	set := s.GetMutualFollowers(userId, otherId)
	return set.Contains(otherId)
}

// GetMutualFollowers returns the subset of given follower IDs that the user is already following back.
func (s *userFollowService) GetMutualFollowers(userId int64, followerIds ...int64) *hashset.Set {
	set := hashset.New()
	list := s.Find(sqls.NewCnd().Eq("user_id", userId).In("other_id", followerIds))
	for _, follow := range list {
		set.Add(follow.OtherID)
	}
	return set
}

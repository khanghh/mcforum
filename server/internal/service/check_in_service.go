package service

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/repository"
	"errors"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"bbs-go/common/dates"
	"bbs-go/pkg/web/params"
	"bbs-go/sqls"
)

var CheckInService = newCheckInService()

func newCheckInService() *checkInService {
	return &checkInService{}
}

type checkInService struct {
	m sync.Mutex
}

func (s *checkInService) Get(id int64) *model.CheckIn {
	return repository.CheckInRepository.Get(sqls.DB(), id)
}

func (s *checkInService) Take(where ...interface{}) *model.CheckIn {
	return repository.CheckInRepository.Take(sqls.DB(), where...)
}

func (s *checkInService) Find(cnd *sqls.Cnd) []model.CheckIn {
	return repository.CheckInRepository.Find(sqls.DB(), cnd)
}

func (s *checkInService) FindOne(cnd *sqls.Cnd) *model.CheckIn {
	return repository.CheckInRepository.FindOne(sqls.DB(), cnd)
}

func (s *checkInService) FindPageByParams(params *params.QueryParams) (list []model.CheckIn, paging *sqls.Paging) {
	return repository.CheckInRepository.FindPageByParams(sqls.DB(), params)
}

func (s *checkInService) FindPageByCnd(cnd *sqls.Cnd) (list []model.CheckIn, paging *sqls.Paging) {
	return repository.CheckInRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *checkInService) Count(cnd *sqls.Cnd) int64 {
	return repository.CheckInRepository.Count(sqls.DB(), cnd)
}

func (s *checkInService) Create(t *model.CheckIn) error {
	return repository.CheckInRepository.Create(sqls.DB(), t)
}

func (s *checkInService) Update(t *model.CheckIn) error {
	return repository.CheckInRepository.Update(sqls.DB(), t)
}

func (s *checkInService) Updates(id int64, columns map[string]interface{}) error {
	return repository.CheckInRepository.Updates(sqls.DB(), id, columns)
}

func (s *checkInService) UpdateColumn(id int64, name string, value interface{}) error {
	return repository.CheckInRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *checkInService) Delete(id int64) {
	repository.CheckInRepository.Delete(sqls.DB(), id)
}

func (s *checkInService) CheckIn(userId int64) error {
	s.m.Lock()
	defer s.m.Unlock()
	var (
		checkIn         = s.GetByUserId(userId)
		dayName         = dates.GetDay(time.Now())
		yesterdayName   = dates.GetDay(time.Now().Add(-time.Hour * 24))
		consecutiveDays = 1
		err             error
	)

	if checkIn != nil && checkIn.LatestDayName == dayName {
		return errors.New("You have already checked in")
	}

	if checkIn != nil && checkIn.LatestDayName == yesterdayName {
		consecutiveDays = checkIn.ConsecutiveDays + 1
	}

	if checkIn == nil {
		err = s.Create(&model.CheckIn{
			Model:           model.Model{},
			UserID:          userId,
			LatestDayName:   dayName,
			ConsecutiveDays: consecutiveDays,
			CreateTime:      dates.NowTimestamp(),
			UpdateTime:      dates.NowTimestamp(),
		})
	} else {
		checkIn.LatestDayName = dayName
		checkIn.ConsecutiveDays = consecutiveDays
		checkIn.UpdateTime = dates.NowTimestamp()
		err = s.Update(checkIn)
	}
	if err == nil {
		// Refresh check-in leaderboard cache
		cache.UserCache.RefreshCheckInRank()
		// Handle check-in points
		config := SysConfigService.GetPointConfig()
		if config.CheckInScore > 0 {
			_ = UserService.IncrScore(userId, config.CheckInScore, constants.EntityCheckIn,
				strconv.FormatInt(userId, 10), "Check in "+strconv.Itoa(dayName))

		} else {
			slog.Warn("Check-in points not configured...")
		}
	}
	return err
}

func (s *checkInService) GetByUserId(userId int64) *model.CheckIn {
	return s.FindOne(sqls.NewCnd().Eq("user_id", userId))
}

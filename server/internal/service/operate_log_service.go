package service

import (
	"bbs-go/internal/model"
	"bbs-go/internal/repository"
	"log/slog"
	"net/http"

	"bbs-go/common/dates"
	"bbs-go/common/utils"
	"bbs-go/pkg/web/params"
	"bbs-go/sqls"
)

var OperateLogService = newOperateLogService()

func newOperateLogService() *operateLogService {
	return &operateLogService{}
}

type operateLogService struct {
}

func (s *operateLogService) Get(id int64) *model.OperateLog {
	return repository.OperateLogRepository.Get(sqls.DB(), id)
}

func (s *operateLogService) Take(where ...interface{}) *model.OperateLog {
	return repository.OperateLogRepository.Take(sqls.DB(), where...)
}

func (s *operateLogService) Find(cnd *sqls.Cnd) []model.OperateLog {
	return repository.OperateLogRepository.Find(sqls.DB(), cnd)
}

func (s *operateLogService) FindOne(cnd *sqls.Cnd) *model.OperateLog {
	return repository.OperateLogRepository.FindOne(sqls.DB(), cnd)
}

func (s *operateLogService) FindPageByParams(params *params.QueryParams) (list []model.OperateLog, paging *sqls.Paging) {
	return repository.OperateLogRepository.FindPageByParams(sqls.DB(), params)
}

func (s *operateLogService) FindPageByCnd(cnd *sqls.Cnd) (list []model.OperateLog, paging *sqls.Paging) {
	return repository.OperateLogRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *operateLogService) Count(cnd *sqls.Cnd) int64 {
	return repository.OperateLogRepository.Count(sqls.DB(), cnd)
}

func (s *operateLogService) Create(t *model.OperateLog) error {
	return repository.OperateLogRepository.Create(sqls.DB(), t)
}

func (s *operateLogService) Update(t *model.OperateLog) error {
	return repository.OperateLogRepository.Update(sqls.DB(), t)
}

func (s *operateLogService) Updates(id int64, columns map[string]interface{}) error {
	return repository.OperateLogRepository.Updates(sqls.DB(), id, columns)
}

func (s *operateLogService) UpdateColumn(id int64, name string, value interface{}) error {
	return repository.OperateLogRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *operateLogService) Delete(id int64) {
	repository.OperateLogRepository.Delete(sqls.DB(), id)
}

func (s *operateLogService) AddOperateLog(userId int64, opType, dataType string, dataId int64,
	description string, r *http.Request) {

	operateLog := &model.OperateLog{
		UserId:      userId,
		OpType:      opType,
		DataType:    dataType,
		DataId:      dataId,
		Description: description,
		CreateTime:  dates.NowTimestamp(),
	}
	if r != nil {
		operateLog.Ip = utils.GetRequestIP(r)
		operateLog.UserAgent = utils.GetUserAgent(r)
		operateLog.Referer = r.Header.Get("Referer")
	}
	if err := repository.OperateLogRepository.Create(sqls.DB(), operateLog); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
}

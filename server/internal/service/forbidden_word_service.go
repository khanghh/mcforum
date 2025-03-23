package service

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/repository"
	"regexp"
	"strings"

	"bbs-go/common/strs"
	"bbs-go/sqls"
	"bbs-go/web/params"
)

var ForbiddenWordService = newForbiddenWordService()

func newForbiddenWordService() *forbiddenWordService {
	return &forbiddenWordService{}
}

type forbiddenWordService struct {
}

func (s *forbiddenWordService) Get(id int64) *model.ForbiddenWord {
	return repository.ForbiddenWordRepository.Get(sqls.DB(), id)
}

func (s *forbiddenWordService) Take(where ...interface{}) *model.ForbiddenWord {
	return repository.ForbiddenWordRepository.Take(sqls.DB(), where...)
}

func (s *forbiddenWordService) Find(cnd *sqls.Cnd) []model.ForbiddenWord {
	return repository.ForbiddenWordRepository.Find(sqls.DB(), cnd)
}

func (s *forbiddenWordService) FindOne(cnd *sqls.Cnd) *model.ForbiddenWord {
	return repository.ForbiddenWordRepository.FindOne(sqls.DB(), cnd)
}

func (s *forbiddenWordService) FindPageByParams(params *params.QueryParams) (list []model.ForbiddenWord, paging *sqls.Paging) {
	return repository.ForbiddenWordRepository.FindPageByParams(sqls.DB(), params)
}

func (s *forbiddenWordService) FindPageByCnd(cnd *sqls.Cnd) (list []model.ForbiddenWord, paging *sqls.Paging) {
	return repository.ForbiddenWordRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *forbiddenWordService) Count(cnd *sqls.Cnd) int64 {
	return repository.ForbiddenWordRepository.Count(sqls.DB(), cnd)
}

func (s *forbiddenWordService) Create(t *model.ForbiddenWord) error {
	if err := repository.ForbiddenWordRepository.Create(sqls.DB(), t); err != nil {
		return err
	}
	cache.ForbiddenWordCache.Invalidate()
	return nil
}

func (s *forbiddenWordService) Update(t *model.ForbiddenWord) error {
	if err := repository.ForbiddenWordRepository.Update(sqls.DB(), t); err != nil {
		return err
	}
	cache.ForbiddenWordCache.Invalidate()
	return nil
}

func (s *forbiddenWordService) Updates(id int64, columns map[string]interface{}) error {
	if err := repository.ForbiddenWordRepository.Updates(sqls.DB(), id, columns); err != nil {
		return err
	}
	cache.ForbiddenWordCache.Invalidate()
	return nil
}

func (s *forbiddenWordService) UpdateColumn(id int64, name string, value interface{}) error {
	if err := repository.ForbiddenWordRepository.UpdateColumn(sqls.DB(), id, name, value); err != nil {
		return err
	}
	cache.ForbiddenWordCache.Invalidate()
	return nil
}

func (s *forbiddenWordService) Delete(id int64) {
	repository.ForbiddenWordRepository.Delete(sqls.DB(), id)
	cache.ForbiddenWordCache.Invalidate()
}

func (s forbiddenWordService) Check(content string) (hitWords []string) {
	if strs.IsBlank(content) {
		return
	}
	words := cache.ForbiddenWordCache.Get()
	if len(words) == 0 {
		return
	}
	for _, word := range words {
		if word.Type == constants.ForbiddenWordTypeWord {
			if strings.Contains(content, word.Word) {
				hitWords = append(hitWords, word.Word)
				break
			}
		} else if word.Type == constants.ForbiddenWordTypeRegex {
			// if matched, _ := regexp.MatchString(word.Word, content); matched {
			// 	hitWords = append(hitWords, word.Word)
			// 	break
			// }
			r, _ := regexp.Compile(word.Word)
			if r != nil {
				hits := r.FindAllString(content, 3)
				if len(hits) > 0 {
					hitWords = append(hitWords, hits...)
					break
				}
			}
		}
	}
	return
}

package search

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/model"
	"bbs-go/pkg/markdown"
	"html"
	"log/slog"
	"math"
	"time"

	"bbs-go/common/dates"
	"bbs-go/common/strs"
	"bbs-go/common/utils"
	"bbs-go/sqls"

	"github.com/blevesearch/bleve/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

var index bleve.Index

func Init(indexPath string) {
	var err error
	if index, err = bleve.Open(indexPath); err != nil {
		if err == bleve.ErrorIndexPathDoesNotExist {
			index = newIndex(indexPath)
		} else {
			slog.Error(err.Error())
		}
	}
}

func NewTopicDoc(topic *model.Topic) *TopicDocument {
	if topic == nil {
		return nil
	}
	doc := &TopicDocument{
		Id:          topic.ID,
		ForumId:     topic.ForumId,
		UserId:      topic.UserID,
		Title:       topic.Title,
		Status:      topic.Status,
		Recommended: topic.Recommended,
		CreateTime:  topic.CreateTime,
	}

	// process content
	content := markdown.ToHTML(topic.Content)
	content = utils.GetHtmlText(content)
	content = html.EscapeString(content)

	doc.Content = content

	// process user
	user := cache.UserCache.Get(topic.UserID)
	if user != nil {
		doc.Nickname = user.Nickname
	}

	return doc
}

// IndexData index data
func UpdateTopicIndex(topic *model.Topic) {
	doc := NewTopicDoc(topic)
	if doc == nil {
		return
	}
	err := index.Index(cast.ToString(topic.ID), doc)
	if err != nil {
		slog.Error(err.Error())
	} else {
		slog.Info("add topic search index", slog.Any("id", topic.ID))
	}
}

func DeleteTopicIndex(id int64) error {
	return index.Delete(cast.ToString(id))
}

// paginated query
func SearchTopic(keyword string, nodeId int64, timeRange, page, limit int) (docs []TopicDocument, paging *sqls.Paging, err error) {
	paging = &sqls.Paging{Page: page, Limit: limit}

	query := bleve.NewBooleanQuery()
	query.AddMust(bleve.NewMatchAllQuery())

	if strs.IsNotBlank(keyword) {
		query.AddMust(bleve.NewMatchQuery(keyword))
	}

	if nodeId != 0 {
		if nodeId == -1 { // recommended
			boolFieldQuery := bleve.NewBoolFieldQuery(true)
			boolFieldQuery.SetField("recommended")
			query.AddMust(boolFieldQuery)
		} else {
			f := float64(nodeId)
			b := true
			nodeIdQuery := bleve.NewNumericRangeInclusiveQuery(&f, &f, &b, &b)
			nodeIdQuery.SetField("nodeId")
			query.AddMust(nodeIdQuery)
		}
	}
	if timeRange != 0 {
		var beginTime int64
		if timeRange == 1 { // within one day
			beginTime = dates.Timestamp(time.Now().Add(-24 * time.Hour))
		} else if timeRange == 2 { // within one week
			beginTime = dates.Timestamp(time.Now().Add(-7 * 24 * time.Hour))
		} else if timeRange == 3 { // within one month
			beginTime = dates.Timestamp(time.Now().AddDate(0, -1, 0))
		} else if timeRange == 4 { // within one year
			beginTime = dates.Timestamp(time.Now().AddDate(-1, 0, 0))
		}

		min := float64(beginTime)
		max := float64(math.MaxInt64)
		createTimeQuery := bleve.NewNumericRangeQuery(&min, &max)
		createTimeQuery.SetField("createTime")
		query.AddMust(createTimeQuery)
	}

	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.From = paging.Offset()
	searchRequest.Size = paging.Limit
	searchRequest.Fields = []string{"*"}
	searchRequest.Highlight = bleve.NewHighlightWithStyle("html")
	searchRequest.Highlight.AddField("title")
	searchRequest.Highlight.AddField("content")

	result, err := index.Search(searchRequest)
	if err != nil {
		slog.Error("Search failed:", slog.Any("err", err))
	}

	for _, hit := range result.Hits {

		storedDoc := make(map[string]interface{})
		for key, field := range hit.Fields {
			storedDoc[key] = field
		}

		for field, fragments := range hit.Fragments {
			if len(fragments) > 0 {
				storedDoc[field] = fragments[0]
			}
		}

		if tagField, ok := storedDoc["tags"]; ok {
			switch v := tagField.(type) {
			case string:
				storedDoc["tags"] = []string{v}
			case []interface{}:
				var tags []string
				for _, tag := range v {
					tags = append(tags, tag.(string))
				}
				storedDoc["tags"] = tags
			}
		}

		var doc TopicDocument
		if err := mapstructure.Decode(storedDoc, &doc); err != nil {
			slog.Error(err.Error())
		}
		docs = append(docs, doc)
	}

	return
}

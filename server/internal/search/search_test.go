package search_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/blevesearch/bleve/v2"
)

type TopicDocument struct {
	Id          int64    `json:"id"`
	NodeId      int64    `json:"nodeId"`
	UserId      int64    `json:"userId"`
	Nickname    string   `json:"nickname"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Tags        []string `json:"tags"`
	Recommended bool     `json:"recommended"`
	Status      int      `json:"status"`
	CreateTime  int64    `json:"createTime"`
}

func TestBleve(t *testing.T) {
	// Open or create index
	index, err := bleve.Open("topic_index")
	if err == bleve.ErrorIndexPathDoesNotExist {
		indexMapping := bleve.NewIndexMapping()
		index, err = bleve.New("topic_index", indexMapping)
		// index, err = bleve.NewUsing("topic_index", indexMapping, scorch.Name, scorch.Name, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

	// Add document to index
	doc := TopicDocument{
		Id:          1,
		NodeId:      1,
		UserId:      1,
		Nickname:    "user1",
		Title:       "Example Title",
		Content:     "This is an example content.",
		Tags:        []string{"example", "test"},
		Recommended: true,
		Status:      1,
		CreateTime:  time.Now().Unix(),
	}
	err = index.Index(fmt.Sprintf("%d", doc.Id), doc)
	if err != nil {
		log.Fatal(err)
	}

	// Execute query
	query := bleve.NewMatchQuery("example") // query example keyword "example"
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"id", "title", "content", "createTime"}
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	// Print query results
	for _, hit := range searchResult.Hits {
		var topicDoc TopicDocument
		err := json.Unmarshal(hit.Fields["title"].(json.RawMessage), &topicDoc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Title: %s, Content: %s\n", topicDoc.Title, topicDoc.Content)
	}
}

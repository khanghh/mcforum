package web

import "bbs-go/sqls"

// PageResult paged result data
type PageResult struct {
	Page    *sqls.Paging `json:"page"`    // Pagination info
	Results interface{}  `json:"results"` // Data
}

// CursorResult cursor paged result data
type CursorResult struct {
	Metadata interface{} `json:"metadata,omitempty"`
	Items    interface{} `json:"items"`   // Data
	Cursor   int64       `json:"cursor"`  // Next page
	HasMore  bool        `json:"hasMore"` // Whether there is more data
}

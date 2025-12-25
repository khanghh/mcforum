package sqls

// Paging paged request data
type Paging struct {
	Page  int   `json:"page"`  // Page number
	Limit int   `json:"limit"` // Number of items per page
	Total int64 `json:"total"` // Total number of data items
}

func (p *Paging) Offset() int {
	offset := 0
	if p.Page > 0 {
		offset = (p.Page - 1) * p.Limit
	}
	return offset
}

func (p *Paging) TotalPage() int {
	if p.Total == 0 || p.Limit == 0 {
		return 0
	}
	totalPage := int(p.Total) / p.Limit
	if int(p.Total)%p.Limit > 0 {
		totalPage = totalPage + 1
	}
	return totalPage
}

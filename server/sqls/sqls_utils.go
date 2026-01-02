package sqls

import (
	"database/sql"

	"bbs-go/common/strs"
)

func SqlNullString(value string) sql.NullString {
	return sql.NullString{
		String: value,
		Valid:  len(value) > 0,
	}
}

func SqlNullInt64(value int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: value,
		Valid: value != 0,
	}
}

func KeywordWrap(keyword string) string {
	if strs.IsBlank(keyword) {
		return keyword
	}
	return "`" + keyword + "`"
}

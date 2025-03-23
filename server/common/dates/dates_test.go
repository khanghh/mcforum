package dates_test

import (
	"fmt"
	"testing"
	"time"

	"bbs-go/common/dates"
)

func TestWithTimeAsEndOfDay(t *testing.T) {
	fmt.Println(dates.Timestamp(dates.WithTimeAsEndOfDay(time.Now())))
}

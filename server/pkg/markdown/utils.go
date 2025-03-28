package markdown

import (
	"sync"

	"bbs-go/common/strs"
	"bbs-go/common/utils"

	"github.com/88250/lute"
)

var (
	engine *lute.Lute
	once   sync.Once
)

func getEngine() *lute.Lute {
	once.Do(func() {
		engine = lute.New(func(lute *lute.Lute) {
			// lute.SetToC(true)
			lute.SetSanitize(true)
			lute.SetGFMTaskListItem(true)
		})
	})
	return engine
}

func ToHTML(markdownStr string) string {
	if strs.IsBlank(markdownStr) {
		return ""
	}
	return getEngine().MarkdownStr("", markdownStr)
}

func GetSummary(markdownStr string, summaryLen int) string {
	htmlStr := ToHTML(markdownStr)
	return utils.GetSummaryHtml(htmlStr, summaryLen)
}

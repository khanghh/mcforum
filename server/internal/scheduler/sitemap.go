package scheduler

import (
	"bbs-go/internal/model/constants"
	"bbs-go/internal/uploader"
	"bbs-go/pkg/bbsurls"
	"bytes"
	"compress/gzip"
	"log/slog"
	"strings"
	"time"

	"bbs-go/common/dates"

	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"

	"bbs-go/internal/config"
	"bbs-go/internal/model"
	"bbs-go/internal/service"
)

const (
	changefreqAlways  = "always"
	changefreqHourly  = "hourly"
	changefreqDaily   = "daily"
	changefreqWeekly  = "weekly"
	changefreqMonthly = "monthly"
	changefreqYearly  = "yearly"
	changefreqNever   = "never"
)

var building = false

// GenerateSiteMap
func GenerateSiteMap() {
	if !config.Instance().IsProd() {
		return
	}
	if building {
		slog.Info("Sitemap in building...")
		return
	}
	building = true
	defer func() {
		building = false
	}()

	sm := stm.NewSitemap(0)
	sm.SetDefaultHost(config.Instance().BaseUrl) // website host
	if uploader.IsEnabledOss() {
		sm.SetSitemapsHost(config.Instance().Uploader.AliyunOss.Host) // upload to Aliyun so host set to Aliyun
	} else {
		sm.SetPublicPath(config.Instance().Uploader.Local.Host)
	}
	sm.SetSitemapsPath("sitemap") // sitemap storage directory
	sm.SetVerbose(false)
	sm.SetPretty(false)
	sm.SetCompress(true)
	sm.SetAdapter(&myAdapter{})
	sm.Create()

	sm.Add(stm.URL{
		{"loc", bbsurls.AbsUrl("/")},
		{"lastmod", time.Now()},
		{"changefreq", changefreqHourly},
	})

	sm.Add(stm.URL{
		{"loc", bbsurls.AbsUrl("/topics")},
		{"lastmod", time.Now()},
		{"changefreq", changefreqHourly},
	})

	// var (
	// 	dateFrom = dates.Timestamp(time.Now().AddDate(0, -1, 0))
	// 	dateTo   = dates.NowTimestamp()
	// )

	service.TopicService.ScanDesc(func(topics []model.Topic) {
		for _, topic := range topics {
			if topic.Status == constants.StatusActive {
				topicUrl := bbsurls.AbsTopicUrl(topic.Slug, topic.ID)
				sm.Add(stm.URL{
					{"loc", topicUrl},
					{"lastmod", dates.FromTimestamp(topic.LastCommentTime)},
					{"changefreq", changefreqDaily},
					{"priority", 0.6},
				})
			}
		}
	})

	service.TagService.Scan(func(tags []model.Tag) {
		for _, tag := range tags {
			tagUrl := bbsurls.AbsTagUrl(tag.Name)
			sm.Add(stm.URL{
				{"loc", tagUrl},
				{"lastmod", time.Now()},
				{"changefreq", changefreqMonthly},
			})
		}
	})

	service.UserService.Scan(func(users []model.User) {
		for _, user := range users {
			sm.Add(stm.URL{
				{"loc", bbsurls.AbsUserUrl(user.Username.String)},
				{"lastmod", time.Now()},
				{"changefreq", changefreqWeekly},
			})
		}
	})

	sm.Finalize().PingSearchEngines()
	// sm.Finalize().PingSearchEngines("http://www.google.cn/webmasters/tools/ping?sitemap=%s")
}

// My Adapter
type myAdapter struct {
}

// Bytes gets written content.
func (adp *myAdapter) Bytes() [][]byte {
	return nil
}

// Write will create sitemap xml file into the file systems.
func (adp *myAdapter) Write(loc *stm.Location, data []byte) {
	if stm.GzipPtn.MatchString(loc.Filename()) { // gzip
		var out []byte
		var in bytes.Buffer
		w := gzip.NewWriter(&in)
		_, _ = w.Write(data)
		_ = w.Close()
		out = in.Bytes()

		// write gzip format data
		adp.ossWrite(loc.PathInPublic(), out)

		// write raw data
		adp.ossWrite(strings.ReplaceAll(loc.PathInPublic(), ".gz", ""), data)
	} else { // non-gzip
		adp.ossWrite(loc.PathInPublic(), data)
	}
}

// OSS write
func (adp *myAdapter) ossWrite(fileKey string, out []byte) {
	if _url, err := uploader.PutObject(fileKey, out, ""); err != nil {
		slog.Error("Upload sitemap error:", slog.Any("err", err))
	} else {
		slog.Info("Upload sitemap:", slog.String("url", _url))
	}
}

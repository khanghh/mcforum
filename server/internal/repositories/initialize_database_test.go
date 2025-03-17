package repositories

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/passwd"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

var roles = []map[string]any{
	{
		"id":     1,
		"type":   constants.RoleTypeSystem,
		"name":   "Owner",
		"code":   constants.RoleOwner,
		"remark": "Owner",
		"status": constants.StatusOK,
	},
	{
		"id":     2,
		"type":   constants.RoleTypeSystem,
		"name":   "Admin",
		"code":   constants.RoleAdmin,
		"remark": "Admin",
		"status": constants.StatusOK,
	},
}

var adminUser = models.User{
	Username:      sqls.SqlNullString("admin"),
	Email:         sqls.SqlNullString("admin@mineviet.com"),
	EmailVerified: true,
	Nickname:      "admin",
	Password:      passwd.EncodePassword("123456"),
	Status:        constants.StatusOK,
	Roles:         "owner",
}

var topicNodes = []models.TopicNode{
	{
		Name:        "default",
		Description: "",
		SortNo:      0,
		Status:      constants.StatusOK,
		CreateTime:  dates.NowTimestamp(),
	},
}

var siteNavs, _ = json.Marshal([]models.ActionLink{
	{
		Title: "Homepage",
		Url:   "/",
	},
	{
		Title: "Topics",
		Url:   "/topics",
	},
	{
		Title: "Articles",
		Url:   "/articles",
	},
})

var scoreConfig, _ = json.Marshal(models.ScoreConfig{
	PostTopicScore:   1,
	PostCommentScore: 1,
	CheckInScore:     1,
})

var sysConfigs = []models.SysConfig{
	{
		Key:         constants.SysConfigSiteTitle,
		Value:       "Forum MineViet Network",
		Name:        "Site title",
		Description: "The title of the forum website displayed in the browser tab and metadata.",
	},
	{
		Key:         constants.SysConfigSiteDescription,
		Value:       "Forum MineViet Network",
		Name:        "Site description",
		Description: "A brief description of the forum for SEO and metadata.",
	},
	{
		Key:         constants.SysConfigSiteKeywords,
		Value:       fmt.Sprintf("[%s]", strings.Join([]string{"mineviet"}, ",")),
		Name:        "Site Keywords",
		Description: "A list of keywords for search engine optimization (SEO).",
	},
	{
		Key:         constants.SysConfigSiteNavs,
		Value:       string(siteNavs),
		Name:        "Navigation Menu",
		Description: "JSON configuration for the site's main navigation menu.",
	},
	{
		Key:         constants.SysConfigDefaultNodeId,
		Value:       "1",
		Name:        "Default Topic Category",
		Description: "The default category assigned to new topics if none is specified.",
	},
	{
		Key:         constants.SysConfigTokenExpireDays,
		Value:       "365",
		Name:        "Token Expiration Period",
		Description: "The number of days before authentication tokens expire.",
	},
	{
		Key:         constants.SysConfigScoreConfig,
		Value:       string(scoreConfig),
		Name:        "Points System Configuration",
		Description: "Defines the point rewards for various user actions, such as posting topics, commenting, and daily check-ins.",
	},
}

var menuItems = []models.Menu{
	{
		Title: "Dashboard",
		Name:  "Dashboard",
		Icon:  "icon-dashboard",
		Path:  "/dashboard",
	},
	{
		Title: "User",
		Name:  "User",
		Icon:  "icon-user",
		Path:  "/uer",
	},
	{
		Title: "Post Management",
		Name:  "",
		Icon:  "icon-file",
		Path:  "",
	},
	{
		Title: "TopicNode",
		Name:  "TopicNode",
		Icon:  "",
		Path:  "/topic/topic-node",
	},
	{
		Title: "Topic",
		Name:  "Topic",
		Icon:  "",
		Path:  "/topic/index",
	},
	{
		Title: "Article",
		Name:  "Article",
		Icon:  "icon-nav",
		Path:  "/article",
	},
	{
		Title: "Forbidden Words",
		Name:  "ForbiddenWords",
		Icon:  "icon-stop",
		Path:  "/forbidden-word",
	},
	{
		Title: "Link",
		Name:  "Link",
		Icon:  "icon-link",
		Path:  "/link",
	},
	{
		Title: "Settings",
		Name:  "Settings",
		Icon:  "icon-settings",
		Path:  "/settings",
	},
	{
		Title: "Permission Management",
		Name:  "",
		Icon:  "icon-lock",
		Path:  "",
	},
	{
		Title: "Role",
		Name:  "Role",
		Icon:  "",
		Path:  "/permission/role",
	},
	{
		Title: "Menu",
		Name:  "Menu",
		Icon:  "",
		Path:  "/permission/menu",
	},
	{
		Title: "Permission",
		Name:  "Permission",
		Icon:  "",
		Path:  "/permission/index",
	},
}

// map from roleId to array of menuItemId
var roleMenuConfig = map[int][]int{
	1: {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
	2: {1, 2, 3, 4, 5, 6, 7, 8, 9},
}

func init() {
	connStr := os.Getenv("DB_URL")
	var err error
	db, err = gorm.Open(mysql.Open(connStr), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
}

func TestCreateUser(t *testing.T) {
	adminUser.Id = 1
	adminUser.CreateTime = dates.NowTimestamp()
	adminUser.UpdateTime = dates.NowTimestamp()
	err := UserRepository.Create(db, &adminUser)
	if err != nil {
		panic(err)
	}
}

func TestCreateRoles(t *testing.T) {
	for idx, item := range roles {
		item["sort_no"] = idx
		item["create_time"] = dates.NowTimestamp()
		item["update_time"] = dates.NowTimestamp()
		tx := db.Model(&models.Role{}).Create(item)
		if tx.Error != nil {
			panic(tx.Error)
		}
	}
}

func TestCreateUserRole(t *testing.T) {
	err := UserRoleRepository.Create(db, &models.UserRole{
		Model:      models.Model{1},
		UserId:     adminUser.Id,
		RoleId:     1,
		CreateTime: dates.NowTimestamp(),
	})
	if err != nil {
		panic(err)
	}
}

func TestCreateTopicNodes(t *testing.T) {
	for idx, item := range topicNodes {
		item.Id = int64(idx + 1)
		item.SortNo = idx
		item.CreateTime = dates.NowTimestamp()
		err := TopicNodeRepository.Create(db, &item)
		if err != nil {
			panic(err)
		}
	}
}

func TestCreateSysConfigs(t *testing.T) {
	for idx, item := range sysConfigs {
		item.Id = int64(idx + 1)
		item.CreateTime = dates.NowTimestamp()
		item.UpdateTime = dates.NowTimestamp()
		err := SysConfigRepository.Create(db, &item)
		if err != nil {
			panic(err)
		}
	}
}

func TestCreateMenu(t *testing.T) {
	for idx, item := range menuItems {
		item.Id = int64(idx + 1)
		item.SortNo = idx
		item.CreateTime = dates.NowTimestamp()
		item.UpdateTime = dates.NowTimestamp()
		err := MenuRepository.Create(db, &item)
		if err != nil {
			panic(err)
		}
	}
}

func TestCreateRoleMenu(t *testing.T) {
	idx := 0
	for roleId, menuIds := range roleMenuConfig {
		for _, menuId := range menuIds {
			item := models.RoleMenu{
				Model:      models.Model{int64(idx + 1)},
				RoleId:     int64(roleId),
				MenuId:     int64(menuId),
				CreateTime: dates.NowTimestamp(),
			}
			err := RoleMenuRepository.Create(db, &item)
			if err != nil {
				panic(err)
			}
			idx++
		}
	}
}

func TestInitializeDatabase(t *testing.T) {
	if err := db.AutoMigrate(models.Models...); nil != err {
		panic(err)
	}

	TestCreateUser(t)
	TestCreateRoles(t)
	TestCreateUserRole(t)
	TestCreateTopicNodes(t)
	TestCreateSysConfigs(t)
	TestCreateMenu(t)
	TestCreateRoleMenu(t)

}

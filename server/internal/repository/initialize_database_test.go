package repository

import (
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"encoding/json"
	"os"
	"testing"

	"bbs-go/common/dates"
	"bbs-go/common/passwd"
	"bbs-go/sqls"

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

var adminUser = model.User{
	Username:      sqls.SqlNullString("admin"),
	Email:         sqls.SqlNullString("admin@mineviet.com"),
	EmailVerified: true,
	Nickname:      "admin",
	Password:      passwd.EncodePassword("123456"),
	Status:        constants.StatusOK,
	Roles:         "owner",
}

var forums = []model.Forum{
	{
		Name:        "Thông báo",
		Slug:        "announcement",
		Description: "",
		SortNo:      0,
		Status:      constants.StatusOK,
		CreateTime:  dates.NowTimestamp(),
	},
	{
		Name:        "Hỗ trợ",
		Slug:        "support",
		Description: "",
		SortNo:      0,
		Status:      constants.StatusOK,
		CreateTime:  dates.NowTimestamp(),
	},
	{
		Name:        "Thảo luận",
		Slug:        "discussion",
		Description: "",
		SortNo:      0,
		Status:      constants.StatusOK,
		CreateTime:  dates.NowTimestamp(),
	},
	{
		Name:        "Nhân sự",
		Slug:        "staff",
		Description: "",
		SortNo:      0,
		Status:      constants.StatusOK,
		CreateTime:  dates.NowTimestamp(),
	},
}

func jsonMarshal(val any) string {
	data, _ := json.Marshal(val)
	return string(data)
}

var siteNavs = jsonMarshal([]model.ActionLink{
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

var sysConfigs = []model.SysConfig{
	{
		Key:         constants.SysConfigSiteTitle,
		Value:       "Forum MineViet Network",
		Name:        "Site Title",
		Description: "Title displayed in the browser tab and metadata.",
	},
	{
		Key:         constants.SysConfigSiteDescription,
		Value:       "Forum MineViet Network",
		Name:        "Site Description",
		Description: "Brief site description for SEO and metadata.",
	},
	{
		Key:         constants.SysConfigSiteKeywords,
		Value:       jsonMarshal([]string{"mineviet"}),
		Name:        "SEO Keywords",
		Description: "List of keywords for search engine optimization.",
	},
	{
		Key:         constants.SysConfigSiteNavs,
		Value:       siteNavs,
		Name:        "Navigation Menu",
		Description: "JSON configuration for the main menu.",
	},
	{
		Key:         constants.SysConfigSiteNotification,
		Value:       "<span>website is under maintenance</span>",
		Name:        "Site Notification",
		Description: "Global message for maintenance or updates.",
	},
	{
		Key:         constants.SysConfigRecommendTags,
		Value:       jsonMarshal([]string{"help"}),
		Name:        "Recommended Tags",
		Description: "Predefined tags for content categorization.",
	},
	{
		Key:         constants.SysConfigUrlRedirect,
		Value:       "true",
		Name:        "Show URL Redirect Page",
		Description: "You need to manually confirm whether to go to the external link before jumping",
	},
	{
		Key:         constants.SysConfigScoreConfig,
		Value:       jsonMarshal(model.ScoreConfig{1, 1, 1}),
		Name:        "Points System",
		Description: "Defines points for posts, comments, and check-ins.",
	},
	{
		Key:         constants.SysConfigDefaultForumId,
		Value:       "1",
		Name:        "Default Category",
		Description: "Default category for new topics.",
	},
	{
		Key:         constants.SysConfigArticlePending,
		Value:       "true",
		Name:        "Require Article Approval",
		Description: "Articles need manual approval before publishing.",
	},
	{
		Key:         constants.SysConfigTopicCaptcha,
		Value:       "true",
		Name:        "Require Captcha for Posts",
		Description: "Users must enter a captcha when posting.",
	},
	{
		Key:         constants.SysConfigTokenExpireDays,
		Value:       "7",
		Name:        "Token Expiration",
		Description: "Days before authentication tokens expire.",
	},
	{
		Key:         constants.SysConfigLoginMethod,
		Value:       jsonMarshal(model.LoginMethod{Password: true, QQ: true, Github: true, Osc: true}),
		Name:        "Login Methods",
		Description: "Available authentication methods.",
	},
	{
		Key:         constants.SysConfigCreateTopicEmailVerified,
		Value:       "true",
		Name:        "Email Required for Topics",
		Description: "Users must verify email before posting topics.",
	},
	{
		Key:         constants.SysConfigCreateArticleEmailVerified,
		Value:       "true",
		Name:        "Email Required for Articles",
		Description: "Users must verify email before posting articles.",
	},
	{
		Key:         constants.SysConfigCreateCommentEmailVerified,
		Value:       "true",
		Name:        "Email Required for Comments",
		Description: "Users must verify email before posting comments.",
	},
	{
		Key:         constants.SysConfigModules,
		Value:       jsonMarshal(model.ModulesConfig{Tweet: true, Topic: true, Article: true}),
		Name:        "Enabled Modules",
		Description: "Defines which site modules are active.",
	},
	{
		Key:         constants.SysConfigEnableHideContent,
		Value:       "true",
		Name:        "Allow Hidden Content",
		Description: "Users can attach hidden content to posts.",
	},
	{
		Key:         constants.SysConfigEmailWhitelist,
		Value:       jsonMarshal([]string{"gmail.com"}),
		Name:        "Allowed Email Domains",
		Description: "Permitted domains for user registration.",
	},
}

var menuItems = []model.Menu{
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
		tx := db.Model(&model.Role{}).Create(item)
		if tx.Error != nil {
			panic(tx.Error)
		}
	}
}

func TestCreateUserRole(t *testing.T) {
	err := UserRoleRepository.Create(db, &model.UserRole{
		Model:      model.Model{1},
		UserId:     adminUser.Id,
		RoleId:     1,
		CreateTime: dates.NowTimestamp(),
	})
	if err != nil {
		panic(err)
	}
}

func TestCreateForums(t *testing.T) {
	for idx, item := range forums {
		item.Id = int64(idx + 1)
		item.SortNo = idx
		item.CreateTime = dates.NowTimestamp()
		err := ForumRepository.Create(db, &item)
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
			item := model.RoleMenu{
				Model:      model.Model{int64(idx + 1)},
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
	if err := db.AutoMigrate(model.Models...); nil != err {
		panic(err)
	}

	TestCreateUser(t)
	TestCreateRoles(t)
	TestCreateUserRole(t)
	TestCreateForums(t)
	TestCreateSysConfigs(t)
	TestCreateMenu(t)
	TestCreateRoleMenu(t)

}

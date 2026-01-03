package main

import (
	"bbs-go/internal/model"
	"bbs-go/internal/model/constants"
	"bbs-go/internal/repository"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"bbs-go/common/dates"
	"bbs-go/common/passwd"
	"bbs-go/sqls"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

var roles = []map[string]any{
	{
		"id":     1,
		"type":   constants.RoleTypeSystem,
		"name":   "Owner",
		"code":   constants.RoleOwner,
		"remark": "",
		"status": constants.StatusActive,
	},
	{
		"id":     2,
		"type":   constants.RoleTypeSystem,
		"name":   "Admin",
		"code":   constants.RoleAdmin,
		"remark": "",
		"status": constants.StatusActive,
	},
	{
		"id":     3,
		"type":   constants.RoleTypeSystem,
		"name":   "Admin",
		"code":   constants.RoleModerator,
		"remark": "",
		"status": constants.StatusActive,
	},
}

var adminUser = model.User{
	Model:         model.Model{ID: 1},
	Username:      sqls.SqlNullString("admin"),
	Email:         sqls.SqlNullString("admin@mineviet.com"),
	EmailVerified: true,
	Nickname:      "admin",
	Password:      passwd.EncodePassword("123456"),
	RoleID:        sqls.SqlNullInt64(1),
	IsActive:      true,
}

var forums = []model.Forum{
	{
		Name:        "Thông báo",
		Slug:        "announcement",
		Description: "Nơi đăng các thông báo quan trọng về game, cập nhật phiên bản, sự kiện và bảo trì.",
		SortNo:      0,
		Status:      constants.StatusActive,
		CreateTime:  dates.NowTimestamp(),
	},
	{
		Name:        "Hỗ trợ",
		Slug:        "support",
		Description: "Khu vực hỗ trợ người chơi gặp lỗi, vấn đề tài khoản hoặc cần hướng dẫn khi chơi game.",
		SortNo:      1,
		Status:      constants.StatusActive,
		CreateTime:  dates.NowTimestamp(),
	},
	{
		Name:        "Thảo luận",
		Slug:        "discussion",
		Description: "Diễn đàn dành cho người chơi thảo luận về các chủ đề liên quan đến game, chia sẻ kinh nghiệm và mẹo chơi.",
		SortNo:      2,
		Status:      constants.StatusActive,
		CreateTime:  dates.NowTimestamp(),
	},
	{
		Name:        "Nhân sự",
		Slug:        "staff",
		Description: " Khu vực riêng tư dành cho nhân viên và quản trị viên thảo luận về công việc nội bộ và quản lý cộng đồng.",
		SortNo:      3,
		Status:      constants.StatusActive,
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

var forumMenuItems = []model.MenuItem{
	{
		Name:    "Thông báo",
		URLPath: "/forums/announcement",
		LogoURL: "/icons/announcement.svg",
	},
	{
		Name:    "Thảo luận",
		URLPath: "/forums/discussion",
		LogoURL: "/icons/discussion.svg",
	},
	{
		Name:    "Hỗ trợ",
		URLPath: "/forums/support",
		LogoURL: "/icons/support.svg",
	},
}

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
		Key: constants.SysConfigScoreConfig,
		Value: jsonMarshal(model.ScoreConfig{
			LikeTopicScore:          1,
			LikeCommentScore:        1,
			PostTopicScore:          5,
			PostCommentScore:        2,
			ReceiveTopicLikeScore:   2,
			ReceiveCommentLikeScore: 1,
			CheckInScore:            3,
			Streak7DaysScore:        5,
			Streak30DaysScore:       15,
			ReceiveFollowScore:      2,
			BoostMultiplier:         1.0,
			DailyMaxScore:           10000,
		}),
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
		Key:         constants.SysConfigEnableHiddenContent,
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
	{
		Key:         constants.SysConfigMenuItems,
		Value:       jsonMarshal(forumMenuItems),
		Name:        "Menu Items",
		Description: "Configuration for site menu items.",
	},
	{
		Key:         constants.SysConfigUserObserveSeconds,
		Value:       "10",
		Name:        "test",
		Description: "test",
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
	db, err = gorm.Open(mysql.Open(connStr))
	if err != nil {
		panic(err)
	}
}

func TestCreateUser(t *testing.T) {
	adminUser.ID = 1
	adminUser.CreateTime = dates.NowTimestamp()
	adminUser.UpdateTime = dates.NowTimestamp()
	err := repository.UserRepository.Create(db, &adminUser)
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
	err := repository.UserRoleRepository.Create(db, &model.UserRole{
		Model:      model.Model{1},
		UserID:     adminUser.ID,
		RoleID:     1,
		CreateTime: dates.NowTimestamp(),
	})
	if err != nil {
		panic(err)
	}
}

func TestCreateForums(t *testing.T) {
	for idx, item := range forums {
		item.ID = int64(idx + 1)
		item.SortNo = idx
		item.CreateTime = dates.NowTimestamp()
		err := repository.ForumRepository.Create(db, &item)
		if err != nil {
			panic(err)
		}
	}
}

func TestCreateSysConfigs(t *testing.T) {
	for idx, item := range sysConfigs {
		item.ID = int64(idx + 1)
		item.CreateTime = dates.NowTimestamp()
		item.UpdateTime = dates.NowTimestamp()
		err := repository.SysConfigRepository.Create(db, &item)
		if err != nil {
			fmt.Println(err)
			// panic(err)
		}
	}
}

var adminMenuItems = []model.Menu{
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

func TestCreateMenu(t *testing.T) {
	for idx, item := range adminMenuItems {
		item.ID = int64(idx + 1)
		item.SortNo = idx
		item.CreateTime = dates.NowTimestamp()
		item.UpdateTime = dates.NowTimestamp()
		err := repository.MenuRepository.Create(db, &item)
		if err != nil {
			fmt.Println(err)
			// panic(err)
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
			err := repository.RoleMenuRepository.Create(db, &item)
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
	fmt.Println("aaa")

	// TestCreateRoles(t)
	// TestCreateUser(t)
	// TestCreateUserRole(t)
	// TestCreateForums(t)
	// TestCreateSysConfigs(t)
	// TestCreateMenu(t)
	// TestCreateRoleMenu(t)
}

func TestAddUser(t *testing.T) {
	user := model.User{
		Username:      sqls.SqlNullString("user01"),
		Email:         sqls.SqlNullString(""),
		EmailVerified: false,
		Nickname:      "",
		Password:      passwd.EncodePassword("123456"),
		IsActive:      true,
		RoleID:        sqls.SqlNullInt64(2),
		CreateTime:    dates.NowTimestamp(),
	}
	err := repository.UserRepository.Create(db, &user)
	if err != nil {
		panic(err)
	}

}

func TestGetUser(t *testing.T) {
	var user model.User
	db.Preload("Role").Where("username = ?", "khang").First(&user)
	fmt.Println(user.Role)
}

func TestGetUserByRole(t *testing.T) {
	var users []model.User
	// adminRole
	// var roles []model.Role
	// roleNames := []string{"admin", "owner"}
	err := db.Model(&model.User{}).Preload("Roles", "code IN (?)", "owner").
		Where("username = ?", "asaa").
		Find(&users).
		Error
	if err != nil {
		panic(err)
	}
}

func TestSetForumMenu(t *testing.T) {
	var menuConfig model.SysConfig
	for _, item := range sysConfigs {
		if item.Key == constants.SysConfigMenuItems {
			menuConfig = item
			break
		}
	}
	err := repository.SysConfigRepository.Update(db, &menuConfig)
	if err != nil {
		panic(err)
	}
}

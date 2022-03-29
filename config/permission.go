package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	PermConfig *viper.Viper
)

func init() {
	PermConfig = viper.New()
	PermConfig.SetConfigName("permission")
	PermConfig.SetConfigType("yaml")
	PermConfig.AddConfigPath(".")
	PermConfig.AddConfigPath("./config")
	PermConfig.AddConfigPath("/etc/maintainman/")
	PermConfig.AddConfigPath("$HOME/.maintainman/")

	PermConfig.SetDefault("permission", map[string]any{
		"user": map[string]any{
			"view":      "查看当前用户",
			"create":    "创建用户",
			"update":    "更新用户",
			"updateall": "更新所有用户",
			"role":      "修改角色",
			"division":  "修改部门",
			"delete":    "删除用户",
			"viewall":   "查看所有用户",
			"login":     "登录",
			"register":  "注册",
			"renew":     "更新Token",
		},
		"role": map[string]any{
			"view":    "查看当前角色",
			"create":  "创建角色",
			"update":  "更新角色",
			"delete":  "删除角色",
			"viewall": "查看所有角色",
		},
		"permission": map[string]any{
			"viewall": "查看所有权限",
		},
		"announce": map[string]any{
			"view":    "查看公告",
			"hit":     "点击公告",
			"create":  "创建公告",
			"update":  "更新公告",
			"delete":  "删除公告",
			"viewall": "查看所有公告",
		},
		"order": map[string]any{
			"view":       "查看我的订单",
			"viewfix":    "查看我维修的订单",
			"create":     "创建订单",
			"cancel":     "取消订单",
			"update":     "更新订单",
			"updateall":  "更新所有订单",
			"assign":     "分配订单",
			"selfassign": "给自己分配订单",
			"release":    "释放订单",
			"reject":     "拒绝订单",
			"report":     "上报订单",
			"hold":       "挂起订单",
			"complete":   "完成订单",
			"defect":     "修改故障分类",
			"urgence":    "修改紧急程度",
			"appraise":   "评分",
			"viewall":    "查看所有订单",
		},
		"comment": map[string]any{
			"view":      "查看我的评论",
			"create":    "创建评论",
			"delete":    "删除评论",
			"viewall":   "查看所有评论",
			"createall": "创建所有评论",
			"deleteall": "删除所有评论",
		},
		"tag": map[string]any{
			"create":  "创建标签",
			"delete":  "删除标签",
			"viewall": "查看所有标签",
		},
		"item": map[string]any{
			"create":  "创建零件",
			"delete":  "删除零件",
			"viewall": "查看所有零件",
			"update":  "更新零件",
			"consume": "消耗零件",
		},
	})

	if err := PermConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("Permission configuration file not found: %v\n", err)
			if err := PermConfig.SafeWriteConfig(); err != nil {
				panic(fmt.Errorf("Failed to write permission configuration file: %v", err))
			}
			fmt.Println("Default permission configuration file created.")
		} else {
			panic(fmt.Errorf("Fatal error reading config file: %v", err))
		}
	}
}

package config

import (
	"github.com/spf13/viper"
)

const PermConfigVersion = "1.2.0"

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
			"view":       "查看当前用户",
			"create":     "创建用户",
			"update":     "更新用户",
			"updateall":  "更新所有用户",
			"delete":     "删除用户",
			"viewall":    "查看所有用户",
			"login":      "登录",
			"register":   "注册",
			"wxlogin":    "微信登录",
			"wxregister": "微信注册",
			"renew":      "更新Token",
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
		"image": map[string]any{
			"upload": "上传图片",
			"view":   "查看图片",
			"custom": "处理图片",
		},
		"division": map[string]any{
			"viewall": "查看所有分组",
			"create":  "创建分组",
			"update":  "更新分组",
			"delete":  "删除分组",
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
			"create": "创建标签",
			"delete": "删除标签",
			"view":   "查看标签",
			"add":    "添加标签",
		},
		"item": map[string]any{
			"create":  "创建零件",
			"delete":  "删除零件",
			"viewall": "查看所有零件",
			"update":  "更新零件",
			"consume": "消耗零件",
		},
	})

	ReadAndUpdateConfig(PermConfig, "permission", PermConfigVersion)
}

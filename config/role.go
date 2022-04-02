package config

import (
	"github.com/spf13/viper"
)

const RoleConfigVersion = "1.0.2"

var (
	RoleConfig *viper.Viper
)

func init() {
	RoleConfig = viper.New()
	RoleConfig.SetConfigName("role")
	RoleConfig.SetConfigType("yaml")
	RoleConfig.AddConfigPath(".")
	RoleConfig.AddConfigPath("./config")
	RoleConfig.AddConfigPath("/etc/maintainman/")
	RoleConfig.AddConfigPath("$HOME/.maintainman/")

	RoleConfig.SetDefault("role", []any{
		map[string]any{
			"name":         "banned",
			"display_name": "封停用户",
			"permissions":  []string{},
			"inheritance":  []string{},
		},
		map[string]any{
			"name":         "guest",
			"display_name": "访客",
			"guest":        true,
			"permissions": []string{
				"user.register",
				"user.login",
				"user.wxlogin",
				"user.wxregister",
			},
			"inheritance": []string{},
		},
		map[string]any{
			"name":         "user",
			"display_name": "普通用户",
			"default":      true,
			"permissions": []string{
				"user.view",
				"user.update",
				"user.renew",
				"role.view",
				"announce.view",
				"announce.hit",
				"order.view",
				"order.create",
				"order.cancel",
				"order.update",
				"order.appraise",
				"order.urgence",
				"order.comment.view",
				"order.comment.create",
				"order.comment.delete",
				"tag.view.1",
				"tag.add.1",
			},
			"inheritance": []string{
				"guest",
			},
		},
		map[string]any{
			"name":         "maintainer",
			"display_name": "维护工",
			"permissions": []string{
				"order.reject",
				"order.report",
				"order.complete",
				"item.consume",
				"item.viewall",
				"tag.view.2",
				"tag.add.2",
			},
			"inheritance": []string{
				"user",
			},
		},
		map[string]any{
			"name":         "super_maintainer",
			"display_name": "维护工（可自行接单）",
			"permissions": []string{
				"order.selfassign",
				"order.viewall",
			},
			"inheritance": []string{
				"maintainer",
			},
		},
		map[string]any{
			"name":         "admin",
			"display_name": "管理员",
			"permissions": []string{
				"user.division",
				"announce.*",
				"order.*",
				"tag.*",
				"item.*",
			},
			"inheritance": []string{
				"maintainer",
			},
		},
		map[string]any{
			"name":         "super_admin",
			"display_name": "超级管理员",
			"permissions": []string{
				"*",
			},
			"inheritance": []string{
				"admin",
			},
		},
	})

	ReadAndUpdateConfig(RoleConfig, "role", RoleConfigVersion)
}

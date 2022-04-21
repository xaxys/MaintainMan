package order

import "github.com/spf13/viper"

var orderConfig = viper.New()

func init() {
	orderConfig.SetDefault("item_can_negative", true)

	orderConfig.SetDefault("appraise.timeout", "72h")
	orderConfig.SetDefault("appraise.purge", "1m")
	orderConfig.SetDefault("appraise.default", 5)

	orderConfig.SetDefault("notify.wechat.status.tmpl", "订阅消息模板id")
	orderConfig.SetDefault("notify.wechat.status.order", "模板中 订单编号 字段名")
	orderConfig.SetDefault("notify.wechat.status.content", "模板中 订单标题 字段名")
	orderConfig.SetDefault("notify.wechat.status.status", "模板中 订单状态 字段名")
	orderConfig.SetDefault("notify.wechat.status.time", "模板中 订单更新时间 字段名")
	orderConfig.SetDefault("notify.wechat.status.other", "模板中 备注 字段名 (用于传递维修工信息)")

	orderConfig.SetDefault("notify.wechat.comment.tmpl", "微信留言消息模板id")
	orderConfig.SetDefault("notify.wechat.comment.title", "模板中 订单标题 字段名")
	orderConfig.SetDefault("notify.wechat.comment.name", "模板中 留言人 字段名")
	orderConfig.SetDefault("notify.wechat.comment.messgae", "模板中 留言内容 字段名")
	orderConfig.SetDefault("notify.wechat.comment.time", "模板中 留言时间 字段名")
}

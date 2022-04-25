package wordcloud

import (
	"github.com/xaxys/maintainman/modules/order"
)

type OrderWord struct {
	OrderID   uint        `gorm:"primarykey; comment:订单ID"`
	Order     order.Order `gorm:"foreignKey:OrderID"`
	Content   string      `gorm:"primarykey; comment:词语"`
	WordClass string      `gorm:"primarykey; comment:词性"`
	Count     uint        `gorm:"not null; index; default:0; comment:词频"`
}

type GlobalWord struct {
	Content   string `gorm:"primarykey; comment:词语"`
	WordClass string `gorm:"primarykey; comment:词性"`
	Count     uint   `gorm:"not null; index; default:0; comment:词频"`
}

type WordJson struct {
	Content   string `json:"content"`
	WordClass string `json:"wordclass"`
	Count     uint   `json:"count"`
}

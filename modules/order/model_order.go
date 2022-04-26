package order

import (
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/modules/user"
)

type Order struct {
	model.BaseModel
	UserID       uint       `gorm:"not null; index:idx_order_user_status,priority:1; comment:用户ID"`
	User         *user.User `gorm:"foreignkey:UserID"`
	Title        string     `gorm:"not null; index; size:191; comment:标题"`
	Content      string     `gorm:"not null; comment:内容"`
	Address      string     `gorm:"not null; comment:地址"`
	ContactName  string     `gorm:"not null; size:191; comment:联系人"`
	ContactPhone string     `gorm:"not null; size:191; comment:联系电话"`
	Status       uint       `gorm:"not null; size:5; default:0; index:idx_order_user_status,priority:2; comment:状态 0:非法 1:待处理 2:已接单 3:已完成 4:上报中 5:挂单 6:已取消 7:已拒绝 8:已评价"`
	StatusList   []*Status  `gorm:"foreignkey:OrderID"`
	AllowComment uint       `gorm:"not null; size:2; default:1; comment:是否允许评论 1:允许 2:不允许"`
	Comments     []*Comment `gorm:"foreignkey:OrderID"`
	ItemLogs     []*ItemLog `gorm:"foreignkey:OrderID"`
	Tags         []*Tag     `gorm:"many2many:order_tags;"`
	Appraisal    uint       `gorm:"not null; size:5; default:0; comment:评价 0:未评价 1-5:已评价"`
}

type CreateOrderRequest struct {
	Title        string `json:"title" validate:"required,lte=191"`
	Content      string `json:"content" validate:"omitempty,lte=65535"`
	Address      string `json:"address" validate:"required,lte=65535"`
	ContactName  string `json:"contact_name" validate:"required,lte=191"`
	ContactPhone string `json:"contact_phone" validate:"required,lte=191"`
	Tags         []uint `json:"tags"` // 若干 Tag 的 ID
}

type UpdateOrderRequest struct {
	Title        string `json:"title" validate:"omitempty,lte=191"`
	Content      string `json:"content" validate:"omitempty,lte=65535"`
	Address      string `json:"address" validate:"omitempty,lte=65535"`
	ContactName  string `json:"contact_name" validate:"omitempty,lte=191"`
	ContactPhone string `json:"contact_phone" validate:"omitempty,lte=191"`
	AddTags      []uint `json:"add_tags"` // 若干需要添加的 Tag 的 ID
	DelTags      []uint `json:"del_tags"` // 若干需要删除的 Tag 的 ID
}

type AllOrderRequest struct {
	Title       string `json:"title"       url:"title" validate:"lte=191"`
	UserID      uint   `json:"user_id"     url:"user_id"`
	Status      uint   `json:"status"      url:"status"`      // 状态 0:非法 1:待处理 2:已接单 3:已完成 4:上报中 5:挂单 6:已取消 7:已拒绝 8:已评价
	Tags        []uint `json:"tags"        url:"tags"`        // 若干 Tag 的 ID
	Disjunctive bool   `json:"disjunctive" url:"disjunctive"` // false: 查询包含所有Tag的订单, true: 查询包含任一Tag的订单
	model.PageParam
}

type UserOrderRequest struct {
	Status      uint   `url:"status"`
	Tags        []uint `url:"tags"`
	Disjunctive bool   `url:"disjunctive"`
	model.PageParam
}

type RepairerOrderRequest struct {
	Status      uint   `url:"status"`
	Current     bool   `url:"current"`
	Tags        []uint `url:"tags"`
	Disjunctive bool   `url:"disjunctive"`
	model.PageParam
}

type OrderJson struct {
	ID           uint           `json:"id"`
	UserID       uint           `json:"user_id"`
	User         *user.UserJson `json:"user,omitempty"`
	Title        string         `json:"title"`
	Content      string         `json:"content"`
	Address      string         `json:"address"`
	ContactName  string         `json:"contact_name"`
	ContactPhone string         `json:"contact_phone"`
	Status       uint           `json:"status"`
	AllowComment bool           `json:"allow_comment"`
	CreatedAt    int64          `json:"created_at"` // unix timestamp in seconds (UTC)
	UpdatedAt    int64          `json:"updated_at"` // unix timestamp in seconds (UTC)
	Appraisal    uint           `json:"appraisal"`
	Tags         []*TagJson     `json:"tags,omitempty"`
	Comments     []*CommentJson `json:"comments,omitempty"`
}

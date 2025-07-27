package commons

type NamedCount struct {
	Title string `json:"title"`
	Count int64  `json:"count"`
}

type BasicUser struct {
	TableBase
	ExternalId       string `json:"-"`                             // OAuth外部ID
	Status           int8   `json:"status"`                        // 账号状态: 0 正常；1 禁用
	NickName         string `gorm:"size:32" json:"nickName"`       // 昵称
	Avatar           string `gorm:"size:512" json:"avatar"`        // 头像地址
	Provider         string `gorm:"size:32" json:"provider"`       // 身份源
	Email            string `gorm:"unique" json:"email,omitempty"` // 邮箱
	StripeCustomerId string `json:"-"`                             // Stripe客户ID
}

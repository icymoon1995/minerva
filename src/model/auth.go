package model

import "time"

type Auth struct {
	Id        int `json:"id" xorm:"pk autoincr"`
	Email     string `json:"email"`
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"updated"`
	Status    int       `json:"status" xorm:"Int"`
	Password  string	`json:"password"`
	Telephone string	`json:"telephone"`
}

// 账号激活
const StatusActive = 1

// 账号未激活
const StatusInactive = 0

// 账号冻结
const StatusFreeze = 2

func (*Auth) TableName() string {
	return "auth"
}

func (auth *Auth) IsActive() bool {
	return auth.Status == StatusActive
}

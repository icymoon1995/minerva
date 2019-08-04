package model

import "time"

type Auth struct {
	Id        int `xorm:"pk autoincr"`
	Email     string
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	Status    int       `xorm:"Int"`
	Password  string
	Telephone string
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

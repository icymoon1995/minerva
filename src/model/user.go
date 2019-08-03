package model

import "time"

type User struct {
	Id        int `xorm:"pk autoincr"`
	Name      string
	Avatar    string
	Gender    int
	Email     string
	Birthday  time.Time `xorm:"DATE"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

func (*User) TableName() string {
	return "users"
}

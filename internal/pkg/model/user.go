package model

import (
	"github.com/liaomars/mini-blog/pkg/auth"
	"gorm.io/gorm"
	"time"
)

type UserM struct {
	ID        int64     `gorm:"column:id;primary_key"`
	Username  string    `gorm:"column:username;not null"`
	Password  string    `gorm:"column:password;not null"`
	Nickname  string    `gorm:"column:nickname"`
	Email     string    `gorm:"column:email"`
	Phone     string    `gorm:"column:phone"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

func (u *UserM) TableName() string {
	return "user"
}

func (u *UserM) BeforeCreate(tx *gorm.DB) (err error) {
	u.Password, err = auth.Encrypt(u.Password)
	if err != nil {
		return err
	}
	return nil
}

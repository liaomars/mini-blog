package store

import (
	"context"
	"github.com/liaomars/mini-blog/internal/pkg/model"
	"gorm.io/gorm"
)

// UserStore 定义user store层的接口
type UserStore interface {
	Create(ctx context.Context, user *model.UserM) error
}

type users struct {
	db *gorm.DB
}

var _ UserStore = (*users)(nil)

func newUsers(db *gorm.DB) *users {
	return &users{db}
}

// Create 往数据库写入数据
func (u *users) Create(ctx context.Context, user *model.UserM) error {
	return u.db.Create(&user).Error
}

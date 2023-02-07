package store

import (
	"context"
	"github.com/liaomars/mini-blog/internal/pkg/model"
	"gorm.io/gorm"
)

// UserStore 定义user store层的接口
type UserStore interface {
	Create(ctx context.Context, user *model.UserM) error
	Get(ctx context.Context, username string) (*model.UserM, error)
	Update(ctx context.Context, user *model.UserM) error
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

// Get 通过用户名查询用户数据
func (u *users) Get(ctx context.Context, username string) (*model.UserM, error) {
	var user model.UserM
	if err := u.db.Where("username=?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *users) Update(ctx context.Context, user *model.UserM) error {
	return u.db.Save(user).Error
}

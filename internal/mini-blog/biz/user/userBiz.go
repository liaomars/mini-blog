package user

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/liaomars/mini-blog/internal/mini-blog/store"
	"github.com/liaomars/mini-blog/internal/pkg/errno"
	"github.com/liaomars/mini-blog/internal/pkg/model"
	v1 "github.com/liaomars/mini-blog/pkg/api/miniblog/v1"
	"github.com/liaomars/mini-blog/pkg/auth"
	"github.com/liaomars/mini-blog/pkg/token"
	"gorm.io/gorm"
	"regexp"
)

type UserBiz interface {
	Create(ctx context.Context, r *v1.CreateUserRequest) error
	Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error)
	ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error
	GetUserInfo(ctx context.Context, username string) (*v1.GetUserResponse, error)
}

type userBiz struct {
	ds store.IStore
}

// 确保 userBiz 实现了 UserBiz 接口.
var _ UserBiz = (*userBiz)(nil)

func New(ds store.IStore) *userBiz {
	return &userBiz{ds: ds}
}

// Create 创建用户
func (b *userBiz) Create(ctx context.Context, r *v1.CreateUserRequest) error {
	var userM model.UserM
	_ = copier.Copy(&userM, r)

	if err := b.ds.Users().Create(ctx, &userM); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'username'", err.Error()); match {
			return errno.ErrUserAlreadyExist
		}
		return err
	}
	return nil
}

// Login 是 UserBiz 接口中 `Login` 方法的实现.
func (b *userBiz) Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error) {
	user, err := b.ds.Users().Get(ctx, r.Username)
	if err != nil {
		return nil, errno.ErrUserNotFound
	}

	// 去数据库检查输入的密码是否正确
	if err := auth.Compare(user.Password, r.Password); err != nil {
		return nil, errno.ErrPasswordIncorrect
	}

	// 生成token
	t, err := token.Sign(r.Username)

	if err != nil {
		return nil, errno.ErrSignToken
	}

	return &v1.LoginResponse{Token: t}, nil
}

// ChangePassword 是 UserBiz 接口中 `ChangePassword` 方法的实现.
func (b *userBiz) ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error {
	user, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		return errno.ErrUserNotFound
	}

	// 检查旧密码是否正确
	if err := auth.Compare(user.Password, r.OldPassword); err != nil {
		return errno.ErrPasswordIncorrect
	}

	user.Password, _ = auth.Encrypt(r.NewPassword)

	// 把数据库的旧密码更新为新密码
	if err := b.ds.Users().Update(ctx, user); err != nil {
		return err
	}

	return nil
}

// GetUserInfo 通过用户名获取用户数据
func (b *userBiz) GetUserInfo(ctx context.Context, username string) (*v1.GetUserResponse, error) {
	userM, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrUserNotFound
		}

		return nil, err
	}

	var r v1.GetUserResponse
	_ = copier.Copy(&r, userM)
	r.CreatedAt = userM.CreatedAt.Format("2006-01-02 15:04:05")
	r.UpdatedAt = userM.UpdatedAt.Format("2006-01-02 15:04:05")

	return &r, nil
}

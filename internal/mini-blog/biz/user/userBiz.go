package user

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/liaomars/mini-blog/internal/mini-blog/store"
	"github.com/liaomars/mini-blog/internal/pkg/errno"
	"github.com/liaomars/mini-blog/internal/pkg/model"
	v1 "github.com/liaomars/mini-blog/pkg/api/miniblog/v1"
	"regexp"
)

type UserBiz interface {
	Create(ctx context.Context, r *v1.CreateUserRequest) error
}

type userBiz struct {
	ds store.IStore
}

// 确保 userBiz 实现了 UserBiz 接口.
var _ UserBiz = (*userBiz)(nil)

func New(ds store.IStore) *userBiz {
	return &userBiz{ds: ds}
}

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

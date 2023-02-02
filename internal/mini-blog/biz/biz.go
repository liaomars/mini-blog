package biz

import (
	"github.com/liaomars/mini-blog/internal/mini-blog/biz/user"
	"github.com/liaomars/mini-blog/internal/mini-blog/store"
)

type IBiz interface {
	Users() user.UserBiz
}

// biz 是IBiz的实现
type biz struct {
	ds store.IStore
}

func NewBiz(ds store.IStore) *biz {
	return &biz{ds: ds}
}

func (b *biz) Users() user.UserBiz {
	return user.New(b.ds)
}

package store

import (
	"gorm.io/gorm"
	"sync"
)

var (
	once sync.Once

	S *datastore
)

// IStore 定义一个接口，让其它对象实现接口的规范
type IStore interface {
	Users() UserStore
}

// datastore 是 IStore 的一个具体实现.
type datastore struct {
	db *gorm.DB
}

// 这种写法是保证datastore实现IStore接口
var _ IStore = (*datastore)(nil)

func NewStore(db *gorm.DB) *datastore {
	once.Do(func() {
		S = &datastore{db: db}
	})
	return S
}

// Users 这个返回UserStore类型的理解是：user模块里的 user类型实现了,
// Userstore接口里的Create方法，所以表示 user类型实现了Userstore接口
func (ds *datastore) Users() UserStore {
	return newUsers(ds.db)
}

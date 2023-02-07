package auth

import (
	casbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	adapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"time"
)

const (
	// casbin 访问控制模型.
	aclModel = `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)`
)

type Authz struct {
	*casbin.SyncedEnforcer
}

// NewAuthz 创建一个casbin完成授权的授权器
func NewAuthz(db *gorm.DB) (*Authz, error) {
	// 初始化casbin mysql数据库适配器
	adapter, err := adapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}
	// 初始化casbin执行器
	m, _ := model.NewModelFromString(aclModel)
	enforcer, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	// 从数据库加载控制策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}
	enforcer.StartAutoLoadPolicy(5 * time.Second)

	a := &Authz{enforcer}
	return a, nil
}

// Authorize 用来进行授权检查.
func (a *Authz) Authorize(sub, obj, act string) (bool, error) {
	return a.Enforce(sub, obj, act)
}

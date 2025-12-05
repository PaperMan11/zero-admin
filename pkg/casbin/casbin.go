package casbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	redisv9 "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func GetDefaultModelText() string {
	return `
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act) && myFunc(r.obj, p.obj)
	`
}

func NewCasbinEnforcer(modelText string, policyAdapter interface{}) (*casbin.SyncedCachedEnforcer, error) {
	if modelText == "" {
		modelText = GetDefaultModelText()
	}
	m, err := model.NewModelFromString(modelText)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewSyncedCachedEnforcer(m, policyAdapter)
	if err != nil {
		return nil, err
	}
	e.AddFunction("myFunc", KeyMatchFn)
	e.EnableAutoSave(true)
	_ = e.LoadPolicy()
	return e, nil
}

func MustNewCasbinEnforcer(modelText string, policyAdapter interface{}) *casbin.SyncedCachedEnforcer {
	e, err := NewCasbinEnforcer(modelText, policyAdapter)
	logx.Must(err)
	return e
}

func NewGormAdapter(dsn string) (interface{}, error) {
	return gormadapter.NewAdapter("mysql", dsn, true)
}

func MustNewGormAdapter(dsn string) interface{} {
	adapter, err := NewGormAdapter(dsn)
	logx.Must(err)
	return adapter
}

func NewRedisWatcher(c *redis.RedisConf, callback func(data string)) (persist.Watcher, error) {
	watcher, err := rediswatcher.NewWatcher(c.Host, rediswatcher.WatcherOptions{
		Options: redisv9.Options{
			Network:  "tcp",
			Password: c.Pass,
		},
		Channel:    "/casbin",
		IgnoreSelf: false,
	})
	if err != nil {
		return nil, err
	}
	if err = watcher.SetUpdateCallback(callback); err != nil {
		return nil, err
	}
	return watcher, nil
}

func MustNewRedisWatcher(c *redis.RedisConf, callback func(data string)) persist.Watcher {
	watcher, err := NewRedisWatcher(c, callback)
	logx.Must(err)
	return watcher
}

func KeyMatchFn(params ...interface{}) (interface{}, error) {
	logx.Debug("------------------ KeyMatchFn ------------")
	name1 := params[0].(string)
	name2 := params[1].(string)
	logx.Debugf("------------------ r.obj: %s", name1)
	logx.Debugf("------------------ p.obj: %s", name2)
	return true, nil
}

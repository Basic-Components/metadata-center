package models

import (
	"context"
	"time"

	proxyerrs "github.com/Basic-Components/connectproxy/errs"
	"github.com/Basic-Components/connectproxy/pgproxy"
	"github.com/delicb/gstring"
	pg "github.com/go-pg/pg/v9"
	orm "github.com/go-pg/pg/v9/orm"

	log "github.com/Basic-Components/components_manager/logger"
	script "github.com/Basic-Components/components_manager/script"
)

// // ErrUserNotMatch 用户不匹配错误
// var ErrUserNotMatch = errors.New("user not match")

// // ErrISSNotMatch 签发机构不匹配错误
// var ErrISSNotMatch = errors.New("iss not match")

// // ErrUnexpectedSigningMethod 使用未知的签名算法
// var ErrUnexpectedSigningMethod = errors.New("Unexpected signing method")

// Component 组件类
type Component struct {
	tableName struct{} `pg:"component"`
	ID        int
	Name      string                 `pg:",notnull"`
	Version   string                 `pg:",notnull"`
	Desc      string                 `pg:"default:''"`
	Services  []*Service             //has many relation
	Schema    map[string]interface{} `pg:",notnull"`
	Image     string                 `pg:"default:''"`
	Ctime     time.Time              `pg:"default:now()"`
	Utime     time.Time              `pg:"default:now()"`
	Dtime     time.Time              `pg:",soft_delete"`
}

/// 自描述

// String 用户的简单描述
func (instance Component) String() string {
	res := gstring.Sprintm(
		"<Component Id={Id:%d};Name={Name};Version={Version}>",
		map[string]interface{}{
			"Id":      instance.ID,
			"Name":    instance.Name,
			"Version": instance.Version,
		})
	return res
}

// Info 组件的简单描述
func (instance *Component) Info() map[string]interface{} {
	res := map[string]interface{}{
		"id":      instance.ID,
		"name":    instance.Name,
		"version": instance.Version,
		"desc":    instance.Desc,
		"image":   instance.Image,
	}
	return res
}

// SchemaInfo 获取组件在网络环境下的schema描述
func (instance *Component) SchemaInfo(hostname string, uri string) map[string]interface{} {
	schema := instance.Schema
	schema["$id"] = hostname + uri
	return schema
}

// InfoWithSchema 服务的的简单描述
func (instance *Component) InfoWithSchema(hostname string, uri string) map[string]interface{} {
	res := instance.Info()
	schema := instance.SchemaInfo(hostname, uri)
	res["schema"] = schema
	return res
}

/// 钩子

// BeforeUpdate 更新数据时记录数据的更新时间
func (instance *Component) BeforeUpdate(ctx context.Context) (context.Context, error) {
	now := time.Now()
	instance.Utime = now
	return ctx, nil
}

//ComponentNewOptions 绑定 JSON 创建新组件的请求结构体
type ComponentNewOptions struct {
	Name    string                 `json:"name" binding:"required"`
	Version string                 `json:"version" binding:"required"`
	Schema  map[string]interface{} `json:"schema" binding:"required"`
	Desc    string                 `json:"desc"`
	Image   string                 `json:"image"`
}

func realcomponentNew(options *ComponentNewOptions) (*Component, error) {
	component := Component{
		Name:    options.Name,
		Version: options.Version,
		Desc:    options.Desc,
		Schema:  options.Schema,
		Image:   options.Image}
	err := pgproxy.DB.Cli.Insert(&component)
	if err != nil {
		return nil, err
	}
	return &component, nil
}

// ComponentNew 创建新的组件
func ComponentNew(options *ComponentNewOptions) (*Component, error) {
	if pgproxy.DB.Ok == false {
		return nil, proxyerrs.ErrProxyNotInited
	}
	return realcomponentNew(options)
}

// ComponentUpdateOptions 更新组件的参数
type ComponentUpdateOptions struct {
	Desc  string `json:"desc"`
	Image string `json:"image"`
}

// Update 更新服务信息
func (instance *Component) Update(options *ComponentUpdateOptions) error {
	if options.Desc != "" {
		instance.Desc = options.Desc
	}
	if options.Image != "" {
		instance.Image = options.Image
	}
	return pgproxy.DB.Cli.Update(instance)
}

// ComponentGetAll 获取全部的组件元信息
func ComponentGetAll() ([]map[string]interface{}, error) {
	if pgproxy.DB.Ok == false {
		return nil, proxyerrs.ErrProxyNotInited
	}
	var components []Component
	res := []map[string]interface{}{}
	err := pgproxy.DB.Cli.Model(&components).Select()
	for _, component := range components {
		temp := component.Info()
		res = append(res, temp)
	}
	return res, err
}

// componentRegistCallback 需要注册到db代理上的回调
func componentRegistCallback(db *pg.DB) error {
	err := db.CreateTable(
		new(Component),
		&orm.CreateTableOptions{FKConstraints: true})
	if err != nil {
		if err.Error() == `ERROR #42P07 relation "component" already exists` {
			return nil
		}
		log.Warn(map[string]interface{}{"error": err, "place": "CreateTable"}, "componentRegistCallback get error")
		return err
	}
	_, err = db.Exec("ALTER TABLE component ADD CONSTRAINT component_name_version_unique unique(name,version)")
	if err != nil {
		return err
	}
	schema := map[string]interface{}{}
	err = json.UnmarshalFromString(script.Schema, &schema)
	if err != nil {
		return err
	}
	_, err = realcomponentNew(&ComponentNewOptions{
		Name:    script.Config.ComponentName,
		Version: script.Config.ComponentVersion,
		Schema:  schema,
	})
	if err != nil {
		return err
	}
	return nil
}

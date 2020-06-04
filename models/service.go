package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	proxyerrs "github.com/Basic-Components/connectproxy/errs"
	"github.com/Basic-Components/connectproxy/etcd3proxy"
	"github.com/Basic-Components/connectproxy/pgproxy"
	"github.com/delicb/gstring"
	pg "github.com/go-pg/pg/v9"
	orm "github.com/go-pg/pg/v9/orm"

	script "github.com/Basic-Components/components_manager/script"

	log "github.com/Basic-Components/components_manager/logger"

	"github.com/xeipuuv/gojsonschema"
)

// Service 服务类
type Service struct {
	tableName   struct{} `pg:"service"`
	ID          int
	Name        string                 `pg:",notnull"`
	Version     string                 `pg:",notnull"`
	Desc        string                 `pg:"default:''"`
	ComponentID int                    `pg:"on_delete:RESTRICT, on_update: CASCADE,notnull"`
	Component   *Component             // has one relation
	Config      map[string]interface{} `pg:",notnull"`
	Ctime       time.Time              `pg:"default:now()"`
	Utime       time.Time              `pg:"default:now()"`
	Dtime       time.Time              `pg:",soft_delete"`
}

/// 自描述

// String 用户的简单描述
func (instance Service) String() string {
	res := gstring.Sprintm(
		"<Service Id={Id:%d};Name={Name};Version={Version};Component={Component}}>",
		map[string]interface{}{
			"Id":        instance.ID,
			"Component": instance.Component.Name + "::" + instance.Component.Version,
			"Name":      instance.Name,
			"Version":   instance.Version,
		})
	return res
}

// Info 服务的的简单描述
func (instance *Service) Info() map[string]interface{} {
	res := map[string]interface{}{
		"id":        instance.ID,
		"name":      instance.Name,
		"version":   instance.Version,
		"desc":      instance.Desc,
		"component": instance.Component.Info(),
	}
	return res
}

// InfoWithConfig 服务的的简单描述
func (instance *Service) InfoWithConfig() map[string]interface{} {
	res := map[string]interface{}{
		"id":        instance.ID,
		"name":      instance.Name,
		"version":   instance.Version,
		"desc":      instance.Desc,
		"component": instance.Component.Info(),
		"config":    instance.Config,
	}
	return res
}

// InfoWithConfigReleaseStatus 服务的的简单描述
func (instance *Service) InfoWithConfigReleaseStatus() map[string]interface{} {
	res := instance.InfoWithConfig()
	online, _ := instance.ConfigReleased(time.Second)
	res["online"] = online
	return res
}

// ConfigString  服务的配置
func (instance *Service) ConfigString() (string, error) {
	config, err := json.MarshalToString(instance.Config)
	if err != nil {
		return "", err
	}
	return config, nil
}

// ReleaseConfigKey 配置发布时的key值
func (instance *Service) ReleaseConfigKey() string {
	return gstring.Sprintm("/confignamespace/{component_name}/{component_version}/{service_name}",
		map[string]interface{}{
			"component_name":    instance.Component.Name,
			"component_version": instance.Component.Version,
			"service_name":      instance.Name,
		})
}

// ReleaseConfigValue 配置发布时的value值
func (instance *Service) ReleaseConfigValue() (string, error) {
	valueMap := map[string]interface{}{
		"meta": map[string]interface{}{
			"component_name":    instance.Component.Name,
			"component_version": instance.Component.Version,
			"service_name":      instance.Name,
			"service_version":   instance.Version,
			"release_time":      time.Now().Unix(),
		},
		"config": instance.Config,
	}
	value, err := json.MarshalToString(valueMap)
	if err != nil {
		return "", err
	}
	return value, nil
}

// ConfigReleased 查看配置是否被发布
func (instance *Service) ConfigReleased(timeout time.Duration) (bool, error) {
	key := instance.ReleaseConfigKey()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := etcd3proxy.Etcd.Cli.Get(ctx, key)
	if err != nil {
		return false, err
	}
	for _, ev := range resp.Kvs {
		if string(ev.Key) == key {
			res := map[string]interface{}{}
			err = json.Unmarshal(ev.Value, &res)
			if err != nil {
				return false, err
			}
			meta, Ok := res["meta"]
			if !Ok {
				return false, errors.New("can not find meta in value")
			}
			serviceversion, Ok := meta.(map[string]interface{})["service_version"]
			if !Ok {
				return false, errors.New("can not find service_version in value")
			}
			if serviceversion.(string) == instance.Version {
				return true, nil
			}
			return false, nil
		}
	}
	return false, errors.New("not find match key")
}

// ReleaseConfig 发布配置
func (instance *Service) ReleaseConfig(timeout time.Duration) error {
	key := instance.ReleaseConfigKey()
	value, err := instance.ReleaseConfigValue()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err = etcd3proxy.Etcd.Cli.Put(ctx, key, value)
	if err != nil {
		return err
	}
	return nil
}

// UnreleaseConfig 下线配置
func (instance *Service) UnreleaseConfig(timeout time.Duration) error {
	key := instance.ReleaseConfigKey()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := etcd3proxy.Etcd.Cli.Get(ctx, key)
	if err != nil {
		return err
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		if string(ev.Key) == key {
			res := map[string]interface{}{}
			err = json.Unmarshal(ev.Value, &res)
			if err != nil {
				return err
			}
			meta, Ok := res["meta"]
			if !Ok {
				return errors.New("can not find meta in value")
			}
			serviceversion, Ok := meta.(map[string]interface{})["service_version"]
			if !Ok {
				return errors.New("can not find service_version in value")
			}
			if serviceversion.(string) == instance.Version {
				_, err := etcd3proxy.Etcd.Cli.Delete(ctx, key)
				if err != nil {
					return err
				}
				return nil
			}
			return errors.New("service version not match")
		}
	}
	return errors.New("not find match key")
}

/// 钩子

// BeforeUpdate 更新数据时记录数据的更新时间
func (instance *Service) BeforeUpdate(ctx context.Context) (context.Context, error) {
	now := time.Now()
	instance.Utime = now
	return ctx, nil
}

// ServiceNewOptions 创建service的参数
type ServiceNewOptions struct {
	Name        string                 `json:"name" binding:"required"`
	Version     string                 `json:"version" binding:"required"`
	ComponentID int                    `json:"component_id" binding:"required"`
	Config      map[string]interface{} `json:"config" binding:"required"`
	Desc        string                 `json:"desc"`
}

func realserviceNew(options *ServiceNewOptions) (*Service, error) {
	componet := Component{ID: options.ComponentID}
	err := pgproxy.DB.Cli.Select(&componet)
	if err != nil {
		return nil, err
	}
	schemaLoader := gojsonschema.NewGoLoader(componet.Schema)
	configLoader := gojsonschema.NewGoLoader(options.Config)
	result, err := gojsonschema.Validate(schemaLoader, configLoader)
	if err != nil {
		log.Info(nil, "Validate error")
		return nil, err
	}
	if !result.Valid() {
		return nil, errors.New("valid wrong")
	}
	service := Service{
		Name:        options.Name,
		Version:     options.Version,
		Desc:        options.Desc,
		ComponentID: componet.ID,
		Component:   &componet,
		Config:      options.Config,
	}
	err = pgproxy.DB.Cli.Insert(&service)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

// ServiceNew 创建新的组件
func ServiceNew(options *ServiceNewOptions) (*Service, error) {
	if pgproxy.DB.Ok == false {
		return nil, proxyerrs.ErrProxyNotInited
	}
	return realserviceNew(options)
}

// ServiceGetAll 获取全部的组件元信息
func ServiceGetAll() ([]map[string]interface{}, error) {
	if pgproxy.DB.Ok == false {
		return nil, proxyerrs.ErrProxyNotInited
	}
	var services []Service
	res := []map[string]interface{}{}
	err := pgproxy.DB.Cli.Model(&services).Column("service.*").Relation("Component").Select()
	for _, service := range services {
		temp := service.Info()
		res = append(res, temp)
	}
	return res, err
}

// ServiceUpdateOptions 更新服务的参数
type ServiceUpdateOptions struct {
	Config map[string]interface{} `json:"config"`
	Desc   string                 `json:"desc"`
}

// Update 更新服务信息
func (instance *Service) Update(options *ServiceUpdateOptions) error {
	if options.Desc != "" {
		instance.Desc = options.Desc
	}
	if options.Config != nil {
		instance.Config = options.Config
	}
	fmt.Println(instance.Config)
	fmt.Println(options.Config)
	return pgproxy.DB.Cli.Update(instance)
}

// serviceRegistCallback 需要注册到db代理上的回调
func serviceRegistCallback(db *pg.DB) error {
	err := db.CreateTable(
		new(Service),
		&orm.CreateTableOptions{FKConstraints: true})
	if err != nil {
		if err.Error() == `ERROR #42P07 relation "service" already exists` {
			return nil
		}
		log.Warn(map[string]interface{}{"error": err, "place": "CreateTable"}, "componentRegistCallback get error")
		return err
	}
	_, err = db.Exec("ALTER TABLE service ADD CONSTRAINT service_name_version_unique unique(name,version)")
	if err != nil {
		log.Warn(map[string]interface{}{"error": err, "place": "CreateUnique constraint"}, "serviceRegistCallback get error")
		return err
	}

	configbytes, err := json.Marshal(script.Config)
	if err != nil {
		return err
	}
	config := map[string]interface{}{}
	err = json.Unmarshal(configbytes, &config)
	if err != nil {
		return err
	}
	_, err = realserviceNew(&ServiceNewOptions{
		Name:        script.Config.ServiceName,
		Version:     "0.0.0",
		ComponentID: 1,
		Config:      config,
	})
	if err != nil {
		return err
	}
	return nil
}

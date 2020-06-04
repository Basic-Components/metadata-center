package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/delicb/gstring"
	"go.etcd.io/etcd/clientv3"
)

// configLoaderCallback etcdv3操作的回调函数
type configLoaderCallback func(cli *clientv3.Client) error

//ConfigLoader 配置加载的代理对象
type ConfigLoader struct {
	Ok          bool
	isWatching  bool
	Options     *clientv3.Config
	Cli         *clientv3.Client
	Key         string
	Prefix      string
	config      map[string]interface{}            //服务的配置
	configLock  sync.RWMutex                      //防止冲突的锁
	configs     map[string]map[string]interface{} //通版本组件下所有服务的配置
	configsLock sync.RWMutex                      //防止冲突的锁
	callBacks   []configLoaderCallback
}

//ConfigWatcherCallback watch操作的回调函数签名
type ConfigWatcherCallback func(proxy *ConfigLoader) error

//New 创建一个loader
func New() *ConfigLoader {
	loader := new(ConfigLoader)
	loader.Ok = false
	loader.isWatching = false
	loader.configLock = sync.RWMutex{}
	loader.configsLock = sync.RWMutex{}
	return loader
}

var Loader = New()

// Close 关闭pg
func (proxy *ConfigLoader) Close() {
	if proxy.Ok {
		proxy.Cli.Close()
	}
}

// setNotWatch 停止watch后设置
func (proxy *ConfigLoader) setNotWatch() {
	if proxy.isWatching {
		proxy.isWatching = false
	}
}

// Init 初始化loader
func (proxy *ConfigLoader) Init(key string, prefix string, cli *clientv3.Client) error {
	if proxy.Ok {
		return errors.New("loader already inited")
	}
	proxy.Key = key
	proxy.Prefix = prefix
	proxy.Cli = cli
	for _, cb := range proxy.callBacks {
		err := cb(proxy.Cli)
		if err != nil {
			log.Println("regist callback get error", err)
		} else {
			log.Println("regist callback done")
		}
	}
	proxy.Ok = true
	return nil
}

// InitFromOptionsOptions InitFromOptions的参数
type InitFromOptionsOptions struct {
	ComponentName    string
	ComponentVersion string
	ServiceName      string
	Options          *clientv3.Config
}

// InitFromOptions 从配置条件初始化代理对象
func (proxy *ConfigLoader) InitFromOptions(options *InitFromOptionsOptions) error {
	cli, err := clientv3.New(*options.Options)
	if err != nil {
		return err
	}
	proxy.Options = options.Options
	key := gstring.Sprintm("/confignamespace/{component_name}/{component_version}/{service_name}",
		map[string]interface{}{
			"component_name":    options.ComponentName,
			"component_version": options.ComponentVersion,
			"service_name":      options.ServiceName,
		})
	prefix := gstring.Sprintm("/confignamespace/{component_name}/{component_version}/",
		map[string]interface{}{
			"component_name":    options.ComponentName,
			"component_version": options.ComponentVersion,
		})
	return proxy.Init(key, prefix, cli)
}

// Regist 注册回调函数,在init执行后执行回调函数
func (proxy *ConfigLoader) Regist(cb configLoaderCallback) {
	proxy.callBacks = append(proxy.callBacks, cb)
}

// GetConfig 获取服务配置
func (proxy *ConfigLoader) GetConfig() map[string]interface{} {
	proxy.configLock.RLock()
	defer proxy.configLock.RUnlock()
	return proxy.config
}

// GetConfigs 获取固定版本组件下各个服务的配置
func (proxy *ConfigLoader) GetConfigs() map[string]map[string]interface{} {
	proxy.configsLock.RLock()
	defer proxy.configsLock.RUnlock()
	return proxy.configs
}

// PullConfig 拉取SERVICE_NAME匹配的配置
func (proxy *ConfigLoader) PullConfig(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := proxy.Cli.Get(ctx, proxy.Key)
	if err != nil {
		return err
	}
	for _, ev := range resp.Kvs {
		if string(ev.Key) == proxy.Key {
			res := map[string]interface{}{}
			err = json.Unmarshal(ev.Value, &res)
			if err != nil {
				return err
			}
			newconfig, Ok := res["config"]
			if !Ok {
				return errors.New("can not find config in value")
			}
			proxy.configLock.Lock()
			defer proxy.configLock.Unlock()
			proxy.config = newconfig.(map[string]interface{})
			return nil
		}
	}
	return errors.New("not find match key")
}

// PullConfigs 拉取组件名和组件版本一致的所有服务配置
func (proxy *ConfigLoader) PullConfigs(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, err := proxy.Cli.Get(ctx, proxy.Prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	if err != nil {
		return err
	}
	result := map[string]map[string]interface{}{}
	for _, ev := range resp.Kvs {
		res := map[string]interface{}{}
		err = json.Unmarshal(ev.Value, &res)
		if err != nil {
			return err
		}
		newconfig, Ok := res["config"]
		if !Ok {
			return errors.New("can not find config in value")
		}
		result[string(ev.Key)] = newconfig.(map[string]interface{})
	}
	proxy.configsLock.Lock()
	defer proxy.configsLock.Unlock()
	proxy.configs = result
	return nil
}

// WatchConfig 监听key,有更新就拉取配置,回调会在有变动后执行
func (proxy *ConfigLoader) WatchConfig(f ConfigWatcherCallback) error {
	if proxy.isWatching {
		return errors.New("loader is alreay watching")
	}
	defer proxy.setNotWatch()
	rch := proxy.Cli.Watch(context.Background(), proxy.Key)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			if ev.Type == clientv3.EventTypePut {
				fmt.Println("get put event")
				if string(ev.Kv.Key) == proxy.Key {
					res := map[string]interface{}{}
					err := json.Unmarshal(ev.Kv.Value, &res)
					if err != nil {
						fmt.Println(err)
					} else {
						newconfig, Ok := res["config"]
						if !Ok {
							fmt.Println("can not find config in value")
						} else {
							proxy.configLock.Lock()
							proxy.config = newconfig.(map[string]interface{})
							proxy.configLock.Unlock()
							err = f(proxy)
							if err != nil {
								fmt.Println("call back error ", err.Error())
							}
						}
					}
				}
			}
		}
	}
	return nil
}

// WatchConfigs 监听prefix,有更新就拉取配置,回调会在有变动后执行
func (proxy *ConfigLoader) WatchConfigs(f ConfigWatcherCallback) error {
	if proxy.isWatching {
		return errors.New("loader is alreay watching")
	}
	defer proxy.setNotWatch()
	rch := proxy.Cli.Watch(context.Background(), proxy.Key)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			switch ev.Type {
			case clientv3.EventTypePut:
				{
					fmt.Println("get put event")
					key := string(ev.Kv.Key)
					if key == proxy.Key {
						res := map[string]interface{}{}
						err := json.Unmarshal(ev.Kv.Value, &res)
						if err != nil {
							fmt.Println(err)
						} else {
							newconfig, Ok := res["config"]
							if !Ok {
								fmt.Println("can not find config in value")
							} else {
								proxy.configsLock.Lock()
								proxy.configs[key] = newconfig.(map[string]interface{})
								proxy.configsLock.Unlock()
								err = f(proxy)
								if err != nil {
									fmt.Println("call back error ", err.Error())
								}
							}
						}
					}
				}
			case clientv3.EventTypeDelete:
				{
					fmt.Println("get delete event")
					key := string(ev.Kv.Key)
					proxy.configsLock.Lock()
					delete(proxy.configs, key)
					proxy.configsLock.Unlock()
					err := f(proxy)
					if err != nil {
						fmt.Println("call back error ", err.Error())
					}
				}
			default:
				{
					fmt.Println("get unknown event ", ev.Type)
				}
			}
		}
	}
	return nil
}

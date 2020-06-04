# 组件配置管理中心

组件配置管理中心

## 特性

1. 提供用于管理组件和服务配置的统一平台
2. 提供sdk用于拉取服务配置

## 概念设计

组件配置管理中心有两个核心概念:

1. 组件
2. 服务

### 组件

组件是用于描述特定用途一类服务的概念,我们可以将其和docker中的镜像类比,或者类似元类的概念,一个组件会包含一个配置的schema,只要是这个组件的服务的配置应该都满足这个schema.


### 服务

服务是组件的实例,类似swarm中service概念,或者类似类的概念,相同的服务都会拉取同一个配置.

### 节点(目前还没用到)

节点是服务的实例,通常服务做负载均衡会需要起多个节点执行程序

## 项目组成

本项目使用agensgraph(相当于pg 10外带图数据库特性)保存组件和服务信息,如果后续需要扩展成元数据管理工具可以利用其图数据的特性构建关系网络,而配置上线分发依靠etcd.

本项目分为2块:

1. 管理后台
2. sdk

### 管理后台

管理后台目前的操作说明可以看后台首页,测试环境在

### sdk

go语言的sdk在文件夹`configloadersdk`中.需要注意etcd3的官方客户端需要在`go.mod`中指定

```txt
replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.3

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0 // indirect
```

才可以打包执行.

sdk提供一个默认全局代理对象`Loader`,它需要使用`Loader.InitFromOptions`来初始化,初始化需要指定自己服务的名字,组件版本,组件名字以及etcd的连接配置,在此之前我们也可以使用`Loader.Regist`注册回调函数,在初始化时执行,初始化之后别忘了`defer Loader.Close()`.

然后我们可以使用如下接口:

+ `GetConfig`读取自己服务的配置
+ `GetConfigs`读取自己服务所属组件及版本的所有服务配置, 这个返回的时一个map的map,第一层的键格式为`/confignamespace/{component_name}/{component_version}/{service_name}`

+ `PullConfig`从etcd拉取自己服务的配置
+ `PullConfigs`从etcd拉取自己服务所属组件及版本的所有服务配置

+ `WatchConfig` 监听etcd上自己服务配置的改动
+ `WatchConfigs` 监听etcd上自己服务所属组件及版本的所有服务配置的改动

需要注意两个`Watch`接口同时只能执行一个.同时对配置的读写都有加锁.

> 例子:


```go
package main //import "business_config_center"

import (
    ...
	"business_config_center/configloadersdk"
	log "business_config_center/logger"
	"go.etcd.io/etcd/clientv3"
)
func main(){
    ....
    // 展示sdk用法
    etcdaddresses := strings.Split(ETCDURL, ",")
	err := configloadersdk.Loader.InitFromOptions(&configloadersdk.InitFromOptionsOptions{
		ComponentName:    ComponentName,
		ComponentVersion: ComponentVersion,
		ServiceName:      ServiceName,
		Options:          &clientv3.Config{Endpoints: etcdaddresses, DialTimeout: 5 * time.Second},
	})
	if err != nil {
		log.Logger.Fatalln("configloadersdk.InitFromOptions error ", err)
	}
	err = configloadersdk.Loader.PullConfig(time.Second)
	if err != nil {
		log.Info(map[string]interface{}{"err": err}, "configloadersdk.Loader.PullConfig error")
	} else {
		log.Info(map[string]interface{}{"loadconfig": configloadersdk.Loader.GetConfig()}, "configloadersdk.Loader.PullConfig done")
	}
	go configloadersdk.Loader.WatchConfig(func(proxy *configloadersdk.ConfigLoader) error {
		log.Info(map[string]interface{}{"loadnewconfig": proxy.GetConfig()}, "configloadersdk.Loader.WatchConfig changes config")
		return nil
	})
	log.Info(nil, "configloadersdk.InitFromOptions done")
}

```

## changelog

### v0.0.0

项目基本特性构建完成
# Service Manager

+ version: 0.0.1
+ status: dev
+ author: hsz
+ email: hsz1287327@gmail.com

## Description

简单的服务管理工具用于管理和监控服务.

keywords:['tools']

## Feature

+ 服务数据模式管理

## 核心概念

使用本工具需要掌握几个核心概念.这样管理服务才不至于因为概念上有偏差而出现歧义

### 服务(service)

服务资源主体

### 服务的标记(tag)

用于区分同名服务资源的标记,一般是版本号,概念参考docker image的tag

### 服务的使用状态(usage_state)

对应标签的服务的使用状态,可以选择的

### 数据模式(schema)

服务的各项任务需要从外界获取数据或者向外界输出数据,用于描述这些数据样式的数据就是数据模式,本工具使用[json schema](https://json-schema.org/)定义数据模式

#### 数据模式服务的任务(task)

数据模式服务于的细分任务,一个任务可以有输入有输出,也可以是一个任务依赖的任务

任务主要分为`query`,`response`,`receive`,`send`4种,对应请求响应模式和发布订阅模式

#### 数据模式的版本(version)

数据模式有版本区分,以`0.0.0.alpha`这样的形式为标准

#### 数据模式服务的发布环境(Publishing environment)

数据模式对应的服务所在的发布环境.原则上应该先检验发布环境再请求数据模式描述数据.

## 项目组成

项目由如下几个部分组成:

+ 服务部分 项目的主体
+ sdk部分 监控管理的服务应该使用sdk与服务交互
+ client 监控管理使用的客户端

## TODO

功能方面:
+ 服务发现和监控
+ 服务配置管理
+ 服务部署管理
+ 用户系统/权限管理系统
  
sdk方面:
+ js使用的sdk
+ python使用的sdk
+ golang使用的sdk

客户端方面:

+ web客户端
+ chrome扩展客户端
+ electron桌面端
+ 移动端

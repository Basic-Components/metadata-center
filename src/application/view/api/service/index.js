import Router from 'koa-router'
import TagSource from './tag'
import List from './list'
import Profile from './profile'


let Source = new Router()
Source.get('/', List.get)//查看服务列表
Source.post('/', List.post)//创建服务资源
Source.delete('/', List.delete)//删除服务资源
Source.get('/:serviceName', Profile.get)//查看服务的标签列表
Source.post('/:serviceName', Profile.post)//创建服务的标签
Source.put('/:serviceName', Profile.put)//修改服务的标签
Source.delete('/:serviceName', Profile.delete)//删除服务的标签
Source.use('/:serviceName/:tag', TagSource.routes(), TagSource.allowedMethods())

export default Source
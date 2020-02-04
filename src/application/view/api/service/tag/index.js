import Router from 'koa-router'
import List from './list'
import SchemaSource from './schema'
let Source = new Router()
Source.get('/', List.get)//获取对应tag的信息
Source.put('/', List.put) //修改tag
Source.delete('/', List.delete) //删除tag对应的信息
Source.use('/schema', SchemaSource.routes(), SchemaSource.allowedMethods())
export default Source
import Router from 'koa-router'
import List from './list'
import Profile from './profile'

let Source = new Router()
Source.get('/', List.get)//获取schema的列表,如果有搜索条件则返回满足搜索条件的列表,如果搜索到的结果唯一则返回对应的schema内容
Source.post('/', List.post)//创建schema
Source.get('/:schemaId', Profile.get)//获取id对应的schema数据
Source.delete('/:schemaId', Profile.delete)//删除对应id的schema
Source.get('/:schemaId/meta', Profile.getMeta)//获取对应id的schema的元数据
Source.post('/:schemaId/validate', Profile.postValidate)//使用对应id的schema验证数据
export default Source
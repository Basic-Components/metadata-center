import Router from 'koa-router'
import List from './list'
import Profile from './profile'

let Source = new Router()
Source.get('/', List.get)
//Source.post('/', SchemaList.post)
Source.get('/:task/:status/:version', Profile.get)
Source.delete('/:task/:status/:version', Profile.delete)
Source.get('/:task/:status/:version/meta', Profile.getMeta)
Source.post('/:task/:status/:version/validate', Profile.postValidate)
export default Source
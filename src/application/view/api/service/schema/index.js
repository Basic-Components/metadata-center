import Router from 'koa-router'
import List from './list'
import Profile from './profile'

let Source = new Router()
Source.get('/', List.get)
//Source.post('/', SchemaList.post)
Source.get('/:schemaId', Profile.get)
//Source.put('/:schemaId', SchemaProfile.put)
//Source.delete('/:schemaId', SchemaProfile.delete)
export default Source
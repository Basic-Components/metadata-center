import Router from 'koa-router'
import schemaSource from './schema'
import List from './list'
import Profile from './profile'


let Source = new Router()
Source.get('/', List.get)
//Source.post('/', ServiceList.post)
Source.get('/:serviceId', Profile.get)
//Source.post('/:service', ServiceProfile.post)
Source.use('/:serviceId/schema', schemaSource.routes(), schemaSource.allowedMethods())

export default Source
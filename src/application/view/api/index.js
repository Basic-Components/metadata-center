import Router from 'koa-router'
import serviceSource from './service'
import Main from './main'
let api = new Router()
api.get('/', Main.get)
api.use("/service", serviceSource.routes(), serviceSource.allowedMethods())
export default api
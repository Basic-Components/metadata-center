/**
 * # 使用
 * 
 * + 定义纯异步方法对象
 * + 将异步方法注册到api
 */
import Router from 'koa-router'
import api from './api'
//import stream from './stream'
let router = new Router()
router.use("/api", api.routes(), api.allowedMethods())
//router.use("/stream", stream.routes(), stream.allowedMethods())


export default router
import { getSchema, getServiceByid } from '../../../utils'
/**
 * 服务的数据模式列表接口
 */
export default {
    async get(ctx) {
        let serviceId = ctx.params.serviceId
        let service = await getServiceByid(serviceId, "name")
        if (service) {
            let self = {
                source: ctx.url,
                description: `提供指定服务指定搜索条件的数据模式资源列表`
            }
            Object.assign(self, ctx.query, { service })
            let where = { service_id }
            if ("task" in ctx.query) {
                Object.assign(where, { task: ctx.query.task })
            }
            if ("status" in ctx.query) {
                Object.assign(where, { status: ctx.query.status })
            }
            if ("version" in ctx.query) {
                Object.assign(where, { version: ctx.query.status })
            }
            let schemas = await getSchema(where, 'task', 'status', 'version')
            let related = schemas.map((schema) => {
                let task = schema.task
                let status = schema.status
                let version = schema.version
                let source = ctx.url.endsWith("/") ? ctx.url + `${task}/${status}/${version}` : ctx.url + `/${task}/${status}/${version}`
                let description = "提供指定服务指定任务指定发布状态的指定版本数据模型文本"
                Object.assign(schema, { source, description, service })
                return schema
            })
            if (schemas) {
                let result = {
                    self,
                    related
                }
                ctx.body = JSON.stringify(result)
            } else {
                ctx.throw(404)
            }
        } else {
            ctx.throw(404)
        }
    },
    //todo
    async postMessage(ctx) {
        ctx.body = JSON.stringify({
            schema_list: "ok"
        })
    }
}


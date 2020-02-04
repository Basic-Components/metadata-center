import { getSchema,getServiceByid } from '../../../utils'
/**
 * 服务的数据模式详情接口
 */
export default {
    async get(ctx) {
        let serviceId = ctx.params.serviceId
        let task = ctx.params.task
        let status = ctx.params.status
        let version = ctx.params.version
        let schemas = await getSchema({
            serviceId,
            task,
            status,
            version
        }, 'schema')
        switch (schemas.length) {
            case 0: {
                ctx.throw(404)
            }
                break;
            case 1: {
                let schema = schemas[0]
                Object.assign(schema,{
                    "$id": ctx.request.href,
                })
                ctx.body = JSON.stringify(schema)
            }
                break;

            default: {
                ctx.throw(500,"more than one result")
            }
                break;
        }
        
    },
    //todo
    async delete(ctx) {
        ctx.throw(404)
    },
    //todo
    async getMeta(ctx) {
        let serviceId = ctx.params.serviceId
        let task = ctx.params.task
        let status = ctx.params.status
        let version = ctx.params.version
        try {
            let schemas = await getSchema({
                serviceId,
                task,
                status,
                version
            }, 'createdAt', 'updatedAt')
            switch (schemas.length) {
                case 0: {
                    ctx.throw(404)
                }
                    break;
                case 1: {
                    ctx.body = JSON.stringify(schema[0])
                }
                    break;

                default: {
                    ctx.throw(500)
                }
                    break;
            }
        } catch (error) {
            ctx.throw(500)
        }
    },
    //todo
    async postValidate(ctx) {
        ctx.throw(404)
    },
}


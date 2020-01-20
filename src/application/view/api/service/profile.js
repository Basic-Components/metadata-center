import { getServiceByid } from '../../utils'
/**
 * 服务详情接口
 */
export default {
    async get(ctx) {
        let service_name = ctx.params.serviceId
        let service = await getServiceByid(service_name)

        if (service) {
            let self = {
                src: ctx.href,
                url: ctx.url,
                description: `提供${service_name}服务资源`
            }
            Object.assign(self, { attribute: service, method: { GET:, PUT:, DELETE:} })
            let related = [
                {
                    "source": ctx.url.endsWith("/") ? ctx.url + 'schema' : ctx.url + '/schema',
                    "description": "获取本服务的数据模型资源",
                    "method": "GET"
                }
            ]
            ctx.body = JSON.stringify({
                self,
                related
            })
        } else {
            ctx.throw(404)
        }
    }

}
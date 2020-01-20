import { getService } from '../../utils'
/**
 * 服务列表接口
 */
export default {
    async get(ctx) {
        let self = {
            src: ctx.href,
            url: ctx.url,
            description: "提供服务资源列表"
        }
        Object.assign(self, ctx.query)
        if ("tags" in ctx.query) {
            Object.assign(where, { tags: ctx.query.tags })
        }
        if ("status" in ctx.query) {
            Object.assign(where, { tags: ctx.query.status })
        } else {
            Object.assign(where, {
                status: {
                    $ne: "deprecated"
                }
            })
        }
        let result = await getService({}, "name", "description", "tags", "status")
        let related = result.map(
            (data) => {
                let name = data.name
                let sourceURL = ctx.url.endsWith("/") ? `${ctx.url}${name}` : `${ctx.url}/${name}`
                let sourceSrc = ctx.href.endsWith("/") ? `${ctx.href}${name}` : `${ctx.href}/${name}`
                let source = Object.assign({
                    url: sourceURL,
                    src: sourceSrc,
                    description: `${name}的服务资源`,
                    method: {
                        GET: `获取${name}的服务资源的属性,子资源和可用方法`,
                        PUT: `更新${name}的服务资源属性`,
                        DELETE: `设置${name}的服务资源过期`
                    }
                })
                return source
            }
        )
        ctx.body = JSON.stringify({
            self,
            related
        })
    }
}

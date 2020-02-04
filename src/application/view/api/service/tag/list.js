import { getService } from '../../utils'
/**
 * 服务列表接口
 */
export default {
    async get(ctx) {
        let self = {
            src: ctx.href,
            url: ctx.url,
            description: "提供指定标记服务资源的标签列表"
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
        ctx.body = JSON.stringify({
            self,
            related
        })
    },
    //todo
    async put(ctx){
        ctx.throw(404)
    },
    //todo
    async delete(ctx){
        ctx.throw(404)
    }
}

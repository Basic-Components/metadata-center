export default {
    async get(ctx) {
        let self = {
            src: ctx.href,
            url: ctx.url,
            description: "提供资源列表"
        }
        let related = [
            {
                source: "server",
                src: ctx.href.endsWith("/") ? ctx.href + "service" : ctx.href + "/service",
                url: ctx.url.endsWith("/") ? ctx.url + "service" : ctx.url + "/service",
                description: "服务资源",
                method: {
                    GET: "获取已有的服务资源",
                    POST: "创建服务资源"
                }
            }
        ]
        ctx.body = JSON.stringify({
            self,
            related
        })
    }
}

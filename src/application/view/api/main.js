export default {
    async get(ctx) {
        let self = {
            source: ctx.url,
            description: "提供资源列表"
        }
        let related =[
            {
                source: ctx.url.endsWith("/")? ctx.url+"service": ctx.url+"/service",
                description: "提供服务资源列表"
            }
        ]
        ctx.body = JSON.stringify({
            self,
            related
        })
    }
}

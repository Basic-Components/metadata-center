export default {
    async get(ctx) {
        let self = {
            "source": ctx.url,
            "description": "提供schema的资源详情"
        }
        ctx.body = JSON.stringify({
            self
        })
    }
}

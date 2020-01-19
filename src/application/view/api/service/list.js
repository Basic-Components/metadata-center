import connection from '../../../model'
export default {
    async get(ctx) {
        let self = {
            source: ctx.url,
            description: "提供服务资源列表"
        }
        let find_par = {
            attributes: ['id', 'name', 'description'],
            order: [['updatedAt', 'DESC']],
        }
        let result = await connection.get_table("Service").findAll(find_par)
        let related = result.map((i) => {
            i = i.dataValues
            let id = i.id
            let source = ctx.url.endsWith("/") ? `${ctx.url}${id}` : `${ctx.url}/${id}`
            let obj = Object.assign(i, {
                source: source
            })
            console.log(obj)
            return obj
        })
        ctx.body = JSON.stringify({
            self,
            related
        })
    }
}

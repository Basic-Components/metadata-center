import connection from '../../../model'
export default {
    async get(ctx) {
        let serviceId = parseInt(ctx.params.serviceId)
        let service = await connection.get_table("Service").findByPk(serviceId)

        let self = Object.assign(service.dataValues, { "source": ctx.url })
        let serviceSchemas = await service.getSchemas()
        let related = serviceSchemas.map((i) => {
            i = i.dataValues
            let id = i.id
            let source = ctx.url.endsWith("/") ? `${ctx.url}schema/${id}` : `${ctx.url}/schema/${id}`
            let task = i.task
            let version = i.version
            let status = i.status
            return {
                id,
                source,
                task,
                version,
                status
            }
        })
    ctx.body = JSON.stringify({
        self,
        related
    })
}
}
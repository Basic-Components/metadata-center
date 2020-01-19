import Router from 'koa-router'

import sse from 'koa-sse-stream'


let stream = new Router()
stream.get('/', sse({
    maxClients: 5000,
    pingInterval: 30000
}), AllServiceStream.get)
stream.get('/:service', sse({
        maxClients: 5000,
        pingInterval: 30000
    }), ServiceStream.get)

export default  stream
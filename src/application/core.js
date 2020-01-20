import Koa from 'koa'
import log from 'pino'
import path from 'path'
import stat from 'koa-static'
import koaBody from 'koa-body'
import Pino from 'koa-pino-logger'
import cors from 'koa2-cors'
import error from 'koa-json-error'
//import { updateMap } from "./utils"

class Application {
  /**
   * 
   * @param {string} name 
   * @param {Map} config 
   */
  constructor(config, name = __filename) {
    this.name = name
    if (config instanceof Map) {
      this.config = config
    } else {
      throw "config must be a Map"
    }
    this.logger = log({ level: this.config.get("LOG_LEVEL") })
    this.app = new Koa()
    let staticPath = path.resolve("./", this.config.get("STATIC_PATH"))
    this.staticPath = staticPath
    let static_middle = stat(staticPath)
    this.app.use(static_middle)

    this.app.use(cors({
      origin: this.config.get("ORIGIN"),
      credentials: true
    }))
    this.app.use(error())
    this.app.use(Pino({ level: this.config.get("LOG_LEVEL") }))
    this.app.use(koaBody())
    this.server = require('http').createServer(this.app.callback())
  }

  use(middle) {
    this.app.use(middle)
  }

  /**
   * 
   * @param {string} host 
   * @param {number} port 
   * @param {boolean} debug 
   */
  run(debug = true, port = null, host = null) {
    let h = host ? host : this.config.get("HOST")
    let p = port ? port : this.config.get("PORT")
    if (debug) {
      //log.setLevel(log.levels.INFO)
      this.logger.level = "debug"
    } else {
      this.logger.level = "error"
    }
    this.logger.info("Static Server running @ http://" + h + ":" + p + "/ for path @" + this.staticPath)
    this.server.listen(p, h)
  }
}
export default Application
import Application from "./core"
import router from "./view"
import connection from "./model"


function init_app(config) {
  let app = new Application(config)
  connection.init_app(app)
  connection.init_connection()
  app.use(router.routes())
  app.use(router.allowedMethods())
  return app
}
export default init_app
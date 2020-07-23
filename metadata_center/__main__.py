from scripts import main
from config import load_conf




def main():
    config = load_conf(args)
    app = init_app(config)
    if app.config.WEBSOCKET:
        from sanic.websocket import WebSocketProtocol
        app.run(
            host=app.config.HOST,
            port=app.config.PORT,
            worker=app.config.WORKER,
            protocol = WebSocketProtocol,
            debug=False, access_log=False
        )
    else:
        app.run(
            host=app.config.HOST,
            port=app.config.PORT,
            worker=app.config.WORKER,
            debug=False, access_log=False
        )


if __name__ == "__main__":
    main()

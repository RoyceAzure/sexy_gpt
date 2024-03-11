import threading
from appserver.server import create_app
from appserver.grpc_server import server
from flask import g
import os
from dotenv import load_dotenv
from chatconfig.config_manager import ConfigManager
app = create_app()


@app.teardown_appcontext
def close_db(exception):
    db = g.pop('db', None)
    if db is not None:
        # 关闭数据库连接
        db.close()
        

def run_flask():
    host = os.getenv("FLASK_HOST", "0.0.0.0")
    port = os.getenv("FLASK_PORT", "8082")
    app.run(host = host, port=port, debug=True)
    
def run_grpc(host:str, port: str):
    server(host, port)

if __name__ == '__main__':
    config = ConfigManager.get_config()
    load_dotenv(override=config.IS_ENV_OVERWRITE)
    thread = threading.Thread(target=run_grpc, kwargs={"host":"[::]", "port" : "9092"})
    thread.daemon = True
    thread.start()
    run_flask()

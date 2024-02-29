from appserver.server import create_app
from flask import g
import os
app = create_app()


@app.teardown_appcontext
def close_db(exception):
    db = g.pop('db', None)
    if db is not None:
        # 关闭数据库连接
        db.close()
        
        
if __name__ == '__main__':
    host = os.getenv("FLASK_HOST", "0.0.0.0")
    port = os.getenv("FLASK_PORT", "5000")
    app.run(host = host, port=port, debug=True)

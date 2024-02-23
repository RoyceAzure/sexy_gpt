from appserver.server import create_app
from flask import g

app = create_app()


@app.teardown_appcontext
def close_db(exception):
    db = g.pop('db', None)
    if db is not None:
        # 关闭数据库连接
        db.close()
        
        
if __name__ == '__main__':
    app.run(debug=True)

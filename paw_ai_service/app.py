from flask import Flask, request, jsonify, g
from api.chat_handler import chat_pb
from dotenv import load_dotenv

app = Flask(__name__)
app.register_blueprint(chat_pb, url_prefix = "/api/v1")
load_dotenv(override=True)


@app.route('/echo', methods=['POST'])
def echo():
    # Retrieve the JSON data sent with the POST request
    data = request.json
    
    # Echo back the received data
    return jsonify(data), 200


@app.teardown_appcontext
def close_db(exception):
    db = g.pop('db', None)
    if db is not None:
        # 关闭数据库连接
        db.close()


if __name__ == '__main__':
    app.run(debug=True)

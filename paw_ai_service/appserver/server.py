from flask import Flask
from api.chat_handler import chat_pb


def create_app():
    app = Flask(__name__)
    app.register_blueprint(chat_pb, url_prefix = "/api/v1")
    return app



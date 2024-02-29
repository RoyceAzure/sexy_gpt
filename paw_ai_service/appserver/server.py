from flask import Flask
from api.chat_handler import chat_pb
from dotenv import load_dotenv
from chatconfig.config_manager import ConfigManager

def create_app():
    app = Flask(__name__)
    app.register_blueprint(chat_pb, url_prefix = "/api/v1")
    config = ConfigManager.get_config()
    load_dotenv(override=config.IS_ENV_OVERWRITE)
    
    return app



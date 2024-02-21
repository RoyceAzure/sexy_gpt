# config_manager.py
from chatconfig.config import Constant, DevelopmentConfig, ProductionConfig
import os


class ConfigManager:
    _config = None
    @staticmethod
    def load_config():
        environment = os.getenv('FLASK_ENV', Constant.ENV_DEV)
        if environment == Constant.ENV_DEV:
            ConfigManager._config = DevelopmentConfig()
        else:
            ConfigManager._config = ProductionConfig()

    @staticmethod
    def get_config():
        if ConfigManager._config is None:
            ConfigManager.load_config()
        return ConfigManager._config

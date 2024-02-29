# config.py
from datetime import timedelta
import os

class Constant:
    ENV_DEV = "development"
    ENV_PROD = "production"
    REDIS_HOST= "REDIS_HOST"
    REDIS_PORT = "REDIS_PORT"


class Config:
    FREE_CHAT_MEMORY_TIME = os.getenv("FREE_CHAT_MEMORY_TIME", 0)
    pass    
    
class DevelopmentConfig(Config):
    # 測試
    FREE_AGENT_EXPIRED_TIME = timedelta(hours=1)
    DEBUG = True

class ProductionConfig(Config):
    # 開發
    FREE_AGENT_EXPIRED_TIME = timedelta(hours=2)
    DEBUG = False

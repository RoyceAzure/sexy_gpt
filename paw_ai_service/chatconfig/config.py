# config.py
from datetime import timedelta


class Constant:
    ENV_DEV = "development"
    ENV_PROD = "production"



class Config:
    pass    
    
class DevelopmentConfig(Config):
    # 測試
    FREE_AGENT_EXPIRED_TIME = timedelta(hours=1)
    DEBUG = True

class ProductionConfig(Config):
    # 開發
    FREE_AGENT_EXPIRED_TIME = timedelta(hours=2)
    DEBUG = False

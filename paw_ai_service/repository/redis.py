import threading
import redis
from flask import g
from functools import wraps
import time
from redis.exceptions import ConnectionError, TimeoutError, RedisError
from repository.base_repo import RepoFactory




def retry_on_failure(max_retries=3, delay=1):
    def decorator(func):
        @wraps(func)
        def wrapper(*args, **kwargs):
            retries = 0
            while retries < max_retries:
                try:
                    return func(*args, **kwargs)
                except (ConnectionError, TimeoutError) as e:
                    print(f"Caught connection error: {e}, retrying in {delay} seconds...")
                    time.sleep(delay)
                    retries += 1
                except RedisError as e:
                    print(f"Redis error: {e}")
                    break  # 如果是其他类型的Redis错误，则不重试
            
            return []  
        return wrapper
    return decorator




class RedisRepo:
    """
        使用裝飾器統一處理錯誤
        若有錯誤則回傳None, 不跳error, 後續正常執行
    """
    def __init__(self, host, port, db, password):
        retry_count = 0
        while retry_count < 3:
            try:
                self.pool = redis.ConnectionPool(host=host, port=port, db=db, password=password)
                self.db = redis.Redis(connection_pool=self.pool)
                break  # 连接成功，退出循环
            except (ConnectionError, TimeoutError) as e:
                print(f"Redis connection error: {e}, retrying...")
                retry_count += 1
                time.sleep(1000)  # 等待一段时间后重试
        if retry_count == 3:
            
            print("Failed to connect to Redis after several attempts.")
            
    @retry_on_failure(max_retries=1, delay=1)
    def save(self, key, value):
        """序列化并存储用户对象"""
        self.db.set(key, value)
        
    @retry_on_failure(max_retries=1, delay=1)
    def load(self, key):
        """根据用户 ID 反序列化用户对象"""
        obj = self.db.get(key)
        if obj:
            return obj
        return None
    
    # @retry_on_failure(max_retries=1, delay=1)
    def append(self, key, value):
        try:
            self.db.rpush(key, value)
        except Exception as e:
            print(f"Failed to append to Redis, error {e}")
            
        
    # @retry_on_failure(max_retries=1, delay=1)    
    def get_range(self, key):
        try:
            return self.db.lrange(key,0,-1)
        except Exception as e:
            print(f"Failed to get_range from Redis, error {e}")
            return []
    
    
    # @retry_on_failure(max_retries=1, delay=1)
    def set_expired(self, key, expried_dur):
        try:
            return self.db.expire(key, expried_dur)
        except Exception as e:
            print(f"Failed to get_range from Redis, error {e}")
    
    
class RedisFactory(RepoFactory):
    __redis = None
    __lock = threading.Lock()
    @staticmethod
    def get_scoped_conn(host, port, db=0, password=None):
        if "resdb" not in g:
            repo= RedisRepo(host, port, db, password)
            g.resdb = repo
        return g.resdb
    @classmethod
    def get_singleton_conn(cls, host, port, db=0, password=None):
        if cls.__redis is None:
            with cls.__lock:
                if cls.__redis is None:
                    cls.__redis = RedisRepo(host, port, db, password)
        return cls.__redis
    



import sqlite3
from flask import g
from abc import ABC, abstractmethod
from sqlalchemy import create_engine, text
# 定義全局資料庫連線
# engine = create_engine('sqlite:///paw.db', echo=True)  # echo=True 會顯示詳細的 SQL 日誌，實際部署時可以關閉



class DBConnFactory(ABC):
    @abstractmethod
    def get_scoped_conn(self, url):
        pass
    @abstractmethod
    def get_singleton_conn(self, url):
        pass


class SQLAlchemyConnFactory(DBConnFactory):
    def __init__(self):
        self.__engine = None
    
    def get_scoped_conn(self, url):
        if "db" not in g:
            engine= create_engine(url)
            g.db = engine.connect()
        return g.db
    
    def get_singleton_conn(self, url):
        if self.__engine is None:
            # engine = create_engine('sqlite:///example.db')
            self.__engine = create_engine(url)
        return self.__engine
    

class Repository:
    @staticmethod
    def get_scoped_db(url):
        if "db" not in g:
            g.db = sqlite3.connect(url)
        return g.db
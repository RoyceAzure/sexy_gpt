import sqlite3
import threading
from flask import g
from sqlalchemy import create_engine

from repository.base_repo import RepoFactory



class SQLAlchemyConnFactory(RepoFactory):
    __engine = None
    __lock = threading.Lock()
    @staticmethod
    def get_scoped_conn(url):
        if "db" not in g:
            engine= create_engine(url)
            g.db = engine.connect()
        return g.db
    
    @classmethod
    def get_singleton_conn(cls, url):
        if cls.__engine is None:
            with cls.__lock:
                if cls.__engine is None:
                    cls.__engine = create_engine(url)
        return cls.__engine
    

class Repository:
    @staticmethod
    def get_scoped_db(url):
        if "db" not in g:
            g.db = sqlite3.connect(url)
        return g.db
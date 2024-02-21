import sqlite3
from flask import g
# 定義全局資料庫連線
# engine = create_engine('sqlite:///paw.db', echo=True)  # echo=True 會顯示詳細的 SQL 日誌，實際部署時可以關閉


class Repository:
    @staticmethod
    def get_scoped_db(url):
        if "db" not in g:
            
            g.db = sqlite3.connect(url)
        return g.db
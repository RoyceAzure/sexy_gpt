import sqlite3
from abc import ABC, abstractmethod
from sqlalchemy import inspect, text
from langchain.tools import Tool, tool
from pydantic.v1 import BaseModel
from typing import List






@tool("get_example_tool", return_direct=False)    
def get_sql_exmples():
    """Get examples of user questersions and how to query via run_sqlite_query function"""
    return """
            {"questersion": "列出所有帕魯的名稱", "query": "SELECT name FROM paw_id_name; "},
            {
                "questersion": "列出可以生出翠葉鼠的所有父母組合",
                "query": "SELECT parent1, parent2 FORM breed WHERE child='翠葉鼠' GROUP BY parent1, parent2;",
            },
            {
                "questersion": "我想要生出 企丸王，有哪些父母的組合?",
                "query": "SELECT parent1, parent2 FORM breed WHERE child='企丸王' GROUP BY parent1, parent2;",
            },
            {
                "questersion": "列出所有棉悠悠可以生出的後代",
                "query": "SELECT child FROM breed WHERE parent1='棉悠悠';",
            },
            {
                "questersion": "我手上有一隻葉泥泥，我要怎生出勾魂魷?",
                "query": "SELECT parent2 FROM breed WHERE (parent1='葉泥泥' OR parent2='葉泥泥') AND child = '勾魂魷';",
            },
            {
                "questersion": "伏特喵 跟 瞅什魔 交配可以得到甚麼?",
                "query": "SELECT child FROM breed WHERE parent1='伏特喵' AND parent2 = '瞅什魔';",
            },
            {
                "questersion": "燎火鹿可以配種出花冠龍嗎?",
                "query": "SELECT * FROM breed WHERE parent1='燎火鹿' AND child='花冠龍';",
            },
            {
                "questersion": "冰帝美露帕的繁殖力是多少?",
                "query": "SELECT * FROM fertility WHERE name = '冰帝美露帕';",
            }
            """




class DB_Tools(ABC):
    @abstractmethod
    def get_db_dools(self): 
        pass

class SQLAlchemyTools(DB_Tools):
    """_summary_ 
         return bunch tools for gpt to connect with SQLAlchemy db
    Args:
        DB_Tools (SQL Alchemy engine): db conn
    """
    def __init__(self, engine):
        self.__engine = engine

    def _list_tables(self):
        inspector = inspect(self.__engine)        
        try:
            table_names = inspector.get_table_names()
            return "\n".join(table_names)
        except Exception as err:
            return f"The following error occurred: {str(err)}"
        
    def _run_sqlite_query(self, query):
            try:  
                with self.__engine.connect() as c:
                    result = c.execute(text(query))
                    return result.fetchall()
            except Exception as err:
                return f"The following error occurred: {str(err)}"
    

class SQLAlchemy_Sqllite_Tools(SQLAlchemyTools):  
    def _describe_tables(self, table_names):
        try:
            with self.__engine.connect() as c:
                tables = ", ".join("'" +  table + "'" for table in table_names)
                result = c.execute(text(f"SELECT sql FROM sqlite_master WHERE type='table' AND name in ({tables});"))
                dbResultStr = "\n".join(row[0] for row in result if row[0] is not None)
                template = (
                    f"{dbResultStr}\n  以上你看到的就是所有table schema，"
                    "'breed': 描述 : 表格為一種名為'帕魯'生物的配種表 "
                    "'fertility': 表格內容是各種怕魯的繁殖力 "
                )
                return template
        except Exception as err:
            return f"The following error occurred: {str(err)}"
        
    def get_db_dools(self): 
        class RunQueryArgsSchema(BaseModel):
            query: str
        class DescribeTablesArgsSchema(BaseModel):
            tables_names: List[str]
            
        run_query_tool = Tool.from_function(
            name="run_sqlite_query",
            description="Run a sqlite query.",
            func=self._run_sqlite_query,
            args_schema=RunQueryArgsSchema
        )
        describe_tables_tool = Tool.from_function(
            name="describe_tables",
            description="Given a list of table names, returns the schema of those tables",
            func=self._describe_tables,
            args_schema=DescribeTablesArgsSchema
        )
        
        return [run_query_tool, describe_tables_tool, get_sql_exmples]   




    
    
class SQLAlchemy_Postgres_Tools(SQLAlchemyTools):  
    def _describe_tables(self, table_names):
        try:
            with self.__engine.connect() as c:
                # 準備表名條件字符串
                tables = ", ".join("'" + table + "'" for table in table_names)
                # PostgreSQL 查詢，獲取表結構
                query = text(f"""
                    SELECT table_name, column_name, data_type 
                    FROM information_schema.columns 
                    WHERE table_name IN ({tables}) 
                    ORDER BY table_name, ordinal_position;
                """)
                result = c.execute(query)
                # 組裝結果字符串
                dbResultStr = ""
                current_table = ""
                for row in result:
                    if row['table_name'] != current_table:
                        dbResultStr += f"\nTable '{row['table_name']}':\n"
                        current_table = row['table_name']
                    dbResultStr += f"  {row['column_name']} ({row['data_type']})\n"
                template = (
                    f"{dbResultStr}\n以上你看到的就是所有table schema，"
                    "'breed': 描述 : 表格為一種名為'帕魯'生物的配種表 "
                    "'fertility': 表格內容是各種怕魯的繁殖力 "
                )
                return template
        except Exception as err:
            return f"The following error occurred: {str(err)}"
        
    def get_db_dools(self): 
        class RunQueryArgsSchema(BaseModel):
            query: str
        class DescribeTablesArgsSchema(BaseModel):
            tables_names: List[str]
            
        run_query_tool = Tool.from_function(
            name="run_sqlite_query",
            description="Run a sqlite query.",
            func=self._run_sqlite_query,
            args_schema=RunQueryArgsSchema
        )
        describe_tables_tool = Tool.from_function(
            name="describe_tables",
            description="Given a list of table names, returns the schema of those tables",
            func=self._describe_tables,
            args_schema=DescribeTablesArgsSchema
        )
        
        return [run_query_tool, describe_tables_tool, get_sql_exmples]   
    
    

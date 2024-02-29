import sqlite3
from abc import ABC, abstractmethod
from sqlalchemy import inspect, text
from langchain.tools import Tool
from pydantic.v1 import BaseModel
from typing import List




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
            
    def _describe_tables(self, table_names):
        try:
            with self.__engine.connect() as c:
                tables = ", ".join("'" +  table + "'" for table in table_names)
                result = c.execute(text(f"SELECT sql FROM sqlite_master WHERE type='table' AND name in ({tables});"))
                dbResultStr = "\n".join(row[0] for row in result if row[0] is not None)
                template = (
                    f"{dbResultStr}\n  以上你看到的就是所有table schema，"
                    "'breed': 描述 : 表格為一種名為'帕魯'生物的配種表 "
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
        
        return [run_query_tool, describe_tables_tool]   
    
    

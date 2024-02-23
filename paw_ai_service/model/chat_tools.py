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
    
    
class DB_ToolsA:
    def __init__(self, conn):
        self.conn = conn
        
    def _list_tables(self):
        c = self.conn.cursor()
        try:
            c.execute("SELECT name FROM sqlite_master WHERE type='table';")
            rows = c.fetchall()
            return "\n".join(row[0] for row in rows if row[0] is not None)
        except sqlite3.OperationalError as err:
            return f"The following error occured: {str(err)}"


    # list_tables_tool = Tool.from_function(
    #     name="list_tables",
    #     description="List all tables in db. no need parameter",
    #     func=list_tables
    # )


    def _run_sqlite_query(self, query):
        c = self.conn.cursor()
        try:
            c.execute(query)
            return c.fetchall()
        except sqlite3.OperationalError as err:
            return f"The following error occured: {str(err)}"
        
    # run_query_tool = Tool.from_function(
    #     name="run_sqlite_query",
    #     description="Run a sqlite query.",
    #     func=run_sqlite_query
    # )


    def _describe_tables(self, table_names):
        """
        因為sqlite無法替table加上描述
        可能需要用額外的方式來描述table
        可能要在新建一個table_schema來描述
        """
        try:
            c = self.conn.cursor()
            tables = ", ".join("'" +  table + "'" for table in table_names)
            rows = c.execute(f"SELECT sql FROM sqlite_master WHERE type='table' AND name in ({tables});")
            dbResultStr =  "\n".join(row[0] for row in rows if row[0] is not None)
            template =(
                f"{dbResultStr}\n  以上你看到的就是所有table schema，"
                "'breed': 描述 : 表格為一種名為'帕魯'生物的配種表 "
            )
            return template
        except sqlite3.OperationalError as err:
            return f"The following error occured: {str(err)}"

    # describe_tables_tool = Tool.from_function(
    #     name="describe_tables",
    #     description="Given a list of table names, returns the schema of those tables.",
    #     func=describe_tables
    # )



    # describe_tables_tool = Tool.from_function(
    #     name="describe_tables",
    #     description="Given a list of table names, returns the schema of those tables",
    #     func=describe_tables,
    #     args_schema=DescribeTablesArgsSchema
    # )
    def Getdbtools(self) -> List: 
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
    
    

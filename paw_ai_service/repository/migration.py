from sqlalchemy import create_engine, Column, Integer, String, Float, Text
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
import sqlite3


# sqlite_conn = sqlite3.connect('paw.db')

# sqlite_cur = sqlite_conn.cursor()


Base = declarative_base()

class Breed(Base):
    __tablename__ = 'breed'
    ID = Column(Float, primary_key = True, autoincrement=True)
    parent1 = Column(String)
    parent2 = Column(String)
    child = Column(String)
    __table_args__ = {'extend_existing': True}
class PawIdName(Base):
    __tablename__ = 'paw_id_name'

    ID = Column(Float, primary_key = True)
    name = Column(String)
    __table_args__ = {'extend_existing': True}
class Fertility(Base):
    __tablename__ = 'fertility'

    name = Column(String, primary_key = True)
    fertility = Column(Integer)
    __table_args__ = {'extend_existing': True}
    
    
# engine = create_engine('postgresql+psycopg2://royce:royce@localhost/paw')
# Base.metadata.create_all(engine)

# Session_postgres = sessionmaker(bind=engine)
# session_postgres = Session_postgres()


# sqlite_cur.execute("SELECT * breed fertility")

# breeds = sqlite_cur.fetchall()


# for row in breeds:
#     new_breed = Breed(parent1=row[0], parent2=row[1], child=row[2])
#     session_postgres.add(new_breed)

# session_postgres.commit()




# sqlite_cur.execute("SELECT * FROM fertility")

# fertilitys = sqlite_cur.fetchall()

# fertilitys

# for fer_row in fertilitys:
#     new_fertility = Fertility(name=fer_row[0], fertility=fer_row[1])
#     session_postgres.add(new_fertility)

# session_postgres.commit()


# #需要先手動變換 postgres Id 欄位type為 float, 因為使用SQLAchemy DDL建立時會變成int
# sqlite_cur.execute("SELECT * FROM paw_id_name")

# paw_id_names = sqlite_cur.fetchall()


# for row in paw_id_names[13:]:
#     new_paw_id_name = PawIdName(ID=row[0], name=row[1])
#     session_postgres.add(new_paw_id_name)

# session_postgres.commit()
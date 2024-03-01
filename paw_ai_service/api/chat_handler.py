import uuid
from flask import Blueprint, jsonify
from flask import request
from repository.database import SQLAlchemyConnFactory
from model.chat_tools import SQLAlchemyTools, SQLAlchemy_Sqllite_Tools
from model.openail_model import ChatModelFactory
from model.agent import ChatArgs, SingleTonAgenFactory
from model.prompt import PromptFactory
from model.schema.question_schema import FreetalkSchema, FreetalkNoMemorySchema
from marshmallow import ValidationError
from chatconfig.config_manager import ConfigManager 

chat_pb = Blueprint("chat_pb", __name__)

@chat_pb.route('/chat', methods=['POST'])
def chat_users():
    """
    無memeory，不紀錄紀錄歷史訊息
    """
    qa_schema = FreetalkNoMemorySchema(many=False)
    
    try:
        req = qa_schema.load(request.json)
    except ValidationError as err:
        return jsonify(err.messages), 400
    
    ques = req['question']
    agent_factory = SingleTonAgenFactory.get_singleton_agent_factory()
    # wrap_agent_excutor = agent_factory.get_agent(session_id)
    # if wrap_agent_excutor is None:
    #     return {
    #         "code" : 403,
    #         "message" : "session不存在，請重新建立session"
    #     }
    
    llm = ChatModelFactory.get_singleton_model()
    prompt = PromptFactory.create_chat_prompt_no_memory()
    
    engine = SQLAlchemyConnFactory.get_singleton_conn("sqlite:///paw.db")
    dbTolls = SQLAlchemy_Sqllite_Tools(engine)
    tools = dbTolls.get_db_dools()
    
    config = ConfigManager.get_config()
    chat_arg = ChatArgs(llm, prompt, None, expired_time = config.FREE_CHAT_MEMORY_TIME ,streaming=False, use_memory=False,  tool = tools)
    
    wrap_agent_excutor = agent_factory.build_agent_excutor(chat_arg)
    
    
    result = wrap_agent_excutor.chat_with_agent(ques)

    response = {
        "code":200,
        "ans" : result["output"]
    }

    return response


@chat_pb.route('/chat/init', methods=['GET'])
def init_chat():
    """   
    需要受控管呼叫次數  不然一值呼叫就一直創建  
    """

    # llm = ChatModelFactory.get_openai_model()
    uuid_str = str(uuid.uuid4())
    # prompt = PromptFactory.create_chat_prompt(memory_id = uuid_str)
    
    # sql_factory = SQLAlchemyConnFactory()
    # engine = sql_factory.get_singleton_conn("sqlite:///paw.db")
    # dbTolls = SQLAlchemyTools(engine)
    # tools = dbTolls.get_db_dools()
    
    
    
    
    
    
    # agent_factory = SingleTonAgenFactory.get_singleton_agent_factory()
    # agent_factory.init_wrap_agent(uuid_str, llm, prompt, tools=tools)
    
    response = {
        "code" : 202,
        "sessionid" : uuid_str
    }
    
    return response




@chat_pb.route('/freechat', methods=['POST'])
def free_chat():
    """
    目前使用redis + sessionId, 儲存對話
    對話會保留 (FREE_CHAT_MEMORY_TIME) 秒
    若連接不到redis等同於沒有歷史紀錄，仍然可以work, 但是agent本身有retry 機制  所以回傳答案會比較慢        
    """
    qa_schema = FreetalkSchema(many=False)
    
    try:
        req = qa_schema.load(request.json)
    except ValidationError as err:
        return jsonify(err.messages), 400
    
    session_id = req['session_id']
    ques = req['question']
    

    agent_factory = SingleTonAgenFactory.get_singleton_agent_factory()
    # wrap_agent_excutor = agent_factory.get_agent(session_id)
    # if wrap_agent_excutor is None:
    #     return {
    #         "code" : 403,
    #         "message" : "session不存在，請重新建立session"
    #     }
    
    llm = ChatModelFactory.get_singleton_model()
    prompt = PromptFactory.create_chat_prompt()
    
    engine = SQLAlchemyConnFactory.get_singleton_conn("sqlite:///paw.db")
    dbTolls = SQLAlchemy_Sqllite_Tools(engine)
    tools = dbTolls.get_db_dools()
    
    config = ConfigManager.get_config()
    chat_arg = ChatArgs(llm, prompt, session_id, expired_time = config.FREE_CHAT_MEMORY_TIME ,streaming=False, use_memory=True,  tool = tools)
    
    wrap_agent_excutor = agent_factory.build_agent_excutor(chat_arg)
    
    
    result = wrap_agent_excutor.chat_with_agent(ques)

    response = {
        "code":200,
        "ans" : result["output"]
    }

    return response
import uuid
from flask import Blueprint, jsonify
from flask import request
from db.database import Repository
from model.chat_tools import DB_Tools
from model.openail_model import ChatModelFactory
from langchain.memory import ConversationBufferMemory
from model.agent import AgentFactory, SingleTonAgenFactory
from model.prompt import PromptFactory
from model.schema.question_schema import FreetalkSchema, AnsSchema
from marshmallow import ValidationError


chat_pb = Blueprint("chat_pb", __name__)

@chat_pb.route('/chat', methods=['POST'])
def chat_users():
    """
    db conn 使用scoped
    llm 關乎到key, 要使用不同key嗎?
    每個使用者需要自己一個session  自己一個memory, 應該說memory就是該session紀錄
            
    """
    qa_schema = FreetalkSchema(many=False)
    an_schema = AnsSchema(many=False)
    
    try:
        req = qa_schema.load(request.json)
    except ValueError as err:
        return {"message": "Invalid input", "errors": err.messages}, 400
    
    ques = req['question']
    conn = Repository.get_scoped_db("paw.db")
    dbTolls = DB_Tools(conn)
    tools = dbTolls.Getdbtools()
    llm = ChatModelFactory.get_singleton_model()
    prompt = PromptFactory.create_chat_prompt()
    memory = ConversationBufferMemory(memory_key="chat_history", return_messages=True)
    
    agent_executor = AgentFactory.create_agent_excutor(llm, prompt, tools, memory)
    result = agent_executor(ques)

    response = {
        "ans" : result["output"]
    }
    
    return response


@chat_pb.route('/chat/init', methods=['GET'])
def init_chat():
    """   
    需要受控管呼叫次數  不然一值呼叫就一直創建  
    """

    llm = ChatModelFactory.get_openai_model()
    uuid_str = str(uuid.uuid4())
    prompt = PromptFactory.create_chat_prompt(memory_id = uuid_str)
    agent_factory = SingleTonAgenFactory.get_singleton_agent_factory()
    agent_factory.init_wrap_agent(uuid_str, llm, prompt)
    
    response = {
        "code" : 202,
        "sessionid" : uuid_str
    }
    
    return response




@chat_pb.route('/freechat', methods=['POST'])
def free_chat():
    """
    db conn 使用scoped
    llm 關乎到key, 要使用不同key嗎?
    每個使用者需要自己一個session  自己一個memory, 應該說memory就是該session紀錄
            
    """
    qa_schema = FreetalkSchema(many=False)
    
    try:
        req = qa_schema.load(request.json)
    except ValidationError as err:
        return jsonify(err.messages), 400
    
    session_id = req['session_id']
    ques = req['question']
    
    conn = Repository.get_scoped_db("paw.db")
    dbTolls = DB_Tools(conn)
    tools = dbTolls.Getdbtools()
    
    
    agent_factory = SingleTonAgenFactory.get_singleton_agent_factory()
    wrap_agent_excutor = agent_factory.get_agent(session_id)
    if wrap_agent_excutor is None:
        return {
            "code" : 403,
            "message" : "session不存在，請重新建立session"
        }
    
    wrap_agent_excutor.set_tools(tools)
    result = wrap_agent_excutor.chat_with_agent(ques)

    response = {
        "code":200,
        "ans" : result["output"]
    }

    return response
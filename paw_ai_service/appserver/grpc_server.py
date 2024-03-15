from concurrent import futures
import uuid
import grpc
from marshmallow import ValidationError
from openai import RateLimitError
from chatconfig.config_manager import ConfigManager
from model.agent import ChatArgs, SingleTonAgenFactory
from model.chat_tools import SQLAlchemy_Sqllite_Tools
from model.openail_model import ChatModelFactory
from model.prompt import PromptFactory
from model.schema.question_schema import FreetalkNoMemorySchema, FreetalkSchema
from repository.database import SQLAlchemyConnFactory
from shared.pb.service_pawai_pb2_grpc import PawAIServiceServicer, add_PawAIServiceServicer_to_server
import shared.pb.service_pawai_pb2 as service_pawai_pb2
from grpc_reflection.v1alpha import reflection
import shared.pb.chat_pb2 as chat__pb2

class PawAIServer(PawAIServiceServicer):
    def FreeChat(self, request, context):

        ques = request.question

        agent_factory = SingleTonAgenFactory.get_singleton_agent_factory()
        
        llm = ChatModelFactory.get_singleton_model()
        prompt = PromptFactory.create_chat_prompt_no_memory()
        
        engine = SQLAlchemyConnFactory.get_singleton_conn("sqlite:///paw.db")
        dbTolls = SQLAlchemy_Sqllite_Tools(engine)
        tools = dbTolls.get_db_dools()
        
        config = ConfigManager.get_config()
        chat_arg = ChatArgs(llm, prompt, None, expired_time = config.FREE_CHAT_MEMORY_TIME ,streaming=False, use_memory=False,  tool = tools)
        
        wrap_agent_excutor = agent_factory.build_agent_excutor(chat_arg)
        
        try:
            result = wrap_agent_excutor.chat_with_agent(ques)
            code = 200
            ans = result["output"]
        except RateLimitError as re:
            code = 429
            ans = "not enough of quato"
        except Exception as e:
            code = e.status_code
            ans = "internal err"
            
        response = chat__pb2.ChatResponse(code = code, ans = ans)
        
        return response

    def Chat(self, request, context):
        
        try:
            session_id_raw = request.session_id
            session_id = str(uuid.UUID(session_id_raw, version=4))
        except ValueError :
            context.abort(grpc.StatusCode.INVALID_ARGUMENT, 'session_id must be a valid UUID')
            
            
        ques = request.question
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
        
        try:
            result = wrap_agent_excutor.chat_with_agent(ques)
            code = 200
            ans = result["output"]
        except RateLimitError as re:
            code = 429
            ans = "not enough of quato"
        except Exception as e:
            code = e.status_code
            ans = "internal err"

        response = chat__pb2.ChatResponse(code = code, ans = ans)

        return response
    
    def InitSession(self, request, context):
        """   
        需要受控管呼叫次數  不然一值呼叫就一直創建  
        """

        uuid_str = str(uuid.uuid4())
        
        # response = {
        #     "code" : 202,
        #     "sessionid" : uuid_str
        # }
        
        
        
        res = chat__pb2.InitSessionResponse(code=202, session_id=uuid_str)
        
        return res

def server(host:str, port: str):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    add_PawAIServiceServicer_to_server(PawAIServer(), server)
    
    # 服务名称列表，用于反射
    SERVICE_NAMES = (
        service_pawai_pb2.DESCRIPTOR.services_by_name['PawAIService'].full_name,
        reflection.SERVICE_NAME,
    )

    # 启用服务器反射
    reflection.enable_server_reflection(SERVICE_NAMES, server)
    
    
    server.add_insecure_port(f'{host}:{port}')
    server.start()
    server.wait_for_termination()
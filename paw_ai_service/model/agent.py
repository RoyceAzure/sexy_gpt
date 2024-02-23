from langchain.agents import create_openai_functions_agent, AgentExecutor
from langchain.memory import ConversationBufferMemory
from datetime import datetime
from chatconfig.config_manager import ConfigManager
import threading

class SingleTonAgenFactory:
    _instance = None
    
    @classmethod
    def get_singleton_agent_factory(cls):
        if cls._instance is None:
            cls._instance = AgentFactory()
        return cls._instance 


class WrapAgent:
    def __init__(self, agent, agent_excutor, memory, cr_time):
        self._agent = agent
        self._agent_excutor = agent_excutor
        self._memory = memory
        self._cr_time = cr_time
        self._up_time = cr_time
        
        
    def set_up_time(self, up_time):
        self._up_time = up_time
        
    def _is_agent_expired(self)-> bool:
        config = ConfigManager.get_config()
        if datetime.utcnow() - self._up_time > config.FREE_AGENT_EXPIRED_TIME:
            return True
        return False
    
    def chat_with_agent(self, msg):
        if self._is_agent_expired():
            self._memory.clear()
        return self._agent_excutor.invoke({"input":msg})
        


class AgentFactory:
    def __init__(self) -> None:
        self.lock = threading.Lock()
        self._objects = {} #管理所有創建過的AgentExecutor

    def init_wrap_agent(self, session_id, llm, prompt, tools = [], memory = None):
        """
        如果是已經存在的session, memory由外部處理
        """
        
        if memory is None:
            memory = ConversationBufferMemory(memory_key="chat_history", return_messages=True)
            
        agent, agent_executor = AgentFactory.create_agent_excutor(llm, prompt, tools, memory)
        wrap_agent = WrapAgent(agent, agent_executor, memory, datetime.utcnow())
        with self.lock:
            self._objects[session_id] = wrap_agent
        

    def get_agent(self, id) -> WrapAgent:
        """
            id (str uuid): session id
        """
        if id is None:
            return None
        
        wrap_agent = self._objects.get(id)
        return wrap_agent

    
    def del_agent(self, id) -> bool:
        """
            需要注意 WrapAgent的引用count，只有當其他地方都沒有引用到該 WrapAgent時  才會正確清除
            需要手動設定屬性為None
        """
        try :
            if id is None:
                return True
            with self.lock:
                old_wrap_agent = self.get_agent(id)
                if old_wrap_agent is None:
                    return True
                old_wrap_agent.agent_excutor = None
                self._objects.pop(id, None)
            return True
        except Exception as e:
            print(f"Failed to delete agent {id}: {e}")
            return False

    
    
    @staticmethod
    def create_agent_excutor(llm, prompt, tools = [], memory = None): 
        """create langchain agent excutor

        Args:
            llm (chat model):
            prompt (string): _description_
            tools (langchain.tools.Tool, optional): _description_. Defaults to None.
            memory (ConversationBufferMemory, optional): _description_. Defaults to None.

        Returns:
            AgentExecutor: _description_
        """
        agent = create_openai_functions_agent(
            llm=llm,
            prompt=prompt,
            tools=tools
        )

        """
        memory應該是用來保存AI Message type的資料

        """
        agent_executor = AgentExecutor(
            agent=agent,
            verbose=True,
            tools=tools,
            memory=memory
        )
        
        return agent, agent_executor
from langchain.schema import BaseChatMessageHistory
from langchain.memory import ConversationBufferMemory
from service.redis_sersvice import (
    get_messages_by_conversation_id,
    add_message_to_conversation
)




class RedisMessageHistory(BaseChatMessageHistory):
    conversation_id: str
    expired_time : int
    
    def __init__(self, conversation_id, expired_time):
        self.conversation_id = conversation_id
        self.expired_time = expired_time
    @property
    def messages(self):
        return get_messages_by_conversation_id(self.conversation_id)
    
    def add_message(self, message):
        return add_message_to_conversation(
            conversation_id=self.conversation_id,
            role=message.type,
            content=message.content,
            expired_time = self.expired_time 
        )

    def clear(self):
        pass

def build_memory(session_id, expired_time):
    return ConversationBufferMemory(
        chat_memory=RedisMessageHistory(
            conversation_id=session_id,
            expired_time= expired_time
        ),
        return_messages=True,
        memory_key="chat_history",
        
        # output_key="answer"
    )
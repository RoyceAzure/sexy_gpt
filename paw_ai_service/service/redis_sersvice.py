import datetime
from repository.redis import RedisFactory
import os
from chatconfig.config import Constant
from langchain.schema.messages import AIMessage, HumanMessage, SystemMessage
from pydantic import BaseModel, Field
import uuid
from datetime import datetime, timezone

redis_host = os.environ.get(Constant.REDIS_HOST, "localhost")
redis_port = os.environ.get(Constant.REDIS_PORT,6379)
# rf = RedisFactory()
r = RedisFactory.get_singleton_conn(redis_host, redis_port)


class Message(BaseModel):
    msg_id: str = Field(..., alias='msg_id')
    cr_time : datetime = Field(..., alias='cr_time')
    role: str = Field(..., alias='role')
    content: str =  Field(..., alias='content')
    conversation_id : str = Field(..., alias='conversation_id')
        
    def as_lc_message(self) -> HumanMessage | AIMessage | SystemMessage:
        """
        for open ai message
        """
        if self.role == "human":
            return HumanMessage(content=self.content)
        elif self.role == "ai":
            return AIMessage(content=self.content)
        elif self.role == "system":
            return SystemMessage(content=self.content)
        else:
            raise Exception(f"Unknown message role: {self.role}")


def get_messages_by_conversation_id(
    conversation_id: str,
) -> AIMessage | HumanMessage | SystemMessage:
    """
    Finds all messages that belong to the given conversation_id

    :param conversation_id: The id of the conversation

    :return: A list of messages
    """
#         messages = (
#             db.session.query(Message)
#             .filter_by(conversation_id=conversation_id)
#             .order_by(Message.created_on.desc())
#         )
#         return [message.as_lc_message() for message in messages]

    messages = r.get_range(conversation_id)
    if messages is None:
        return None

    res = []
    for message in messages:
        load_msg = Message.parse_raw(message)
        res.append(load_msg.as_lc_message())
    return res


def add_message_to_conversation(
    conversation_id: str, role: str, content: str, expired_time: int = 0
):
    """
    Creates and stores a new message tied to the given conversation_id
        with the provided role and content

    :param conversation_id: The id of the conversation
    :param role: The role of the message
    :param content: The content of the message

    :return: The created message
    """

    msg_id = str(uuid.uuid4())

    msg = Message(msg_id = msg_id, role = role, 
                  conversation_id = conversation_id, 
                  content= content, 
                  cr_time =  datetime.now(timezone.utc))


    r.append(conversation_id, msg.json(by_alias=False))
    
    if expired_time != 0:
        r.set_expired(conversation_id, expired_time)

    #要測試
#     def get_conversation_components(self,conversation_id: str) -> Dict[str, str]:
#         """
#         Returns the components used in a conversation
#         """
#         conversation = Conversation.find_by(id=conversation_id)
#         return {
#             "llm": conversation.llm,
#             "retriever": conversation.retriever,
#             "memory": conversation.memory,
#         }

    # 意思是llm  retriever  memory都可以直接塞到DB????

#     def set_conversation_components(self,
#         conversation_id: str, llm: str, retriever: str, memory: str
#     ) -> None:
#         """
#         Sets the components used by a conversation
#         """
#         conversation = Conversation.find_by(id=conversation_id)
#         conversation.update(llm=llm, retriever=retriever, memory=memory)
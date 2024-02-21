from dotenv import load_dotenv
from langchain_openai import ChatOpenAI



# memory = ConversationBufferMemory(memory_key="chat_history", return_messages=True)

class ChatModelFactory:
    __model = None
    
    @staticmethod
    def get_singleton_model() -> ChatOpenAI:
        """
            Get a singleton instance of the ChatOpenAI model.
        """
        if ChatModelFactory.__model is None:
            ChatModelFactory.__model = ChatOpenAI()
        return ChatModelFactory.__model
        
    @staticmethod
    def get_openai_model(key=None)-> ChatOpenAI:
        """Get a new instance of the ChatOpenAI model."""
        if key is None:
            return ChatOpenAI()    
        return ChatOpenAI(openai_api_key=key)
        
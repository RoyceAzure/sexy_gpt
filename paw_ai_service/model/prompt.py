from langchain.prompts import (
    ChatPromptTemplate,
    FewShotChatMessagePromptTemplate,
    HumanMessagePromptTemplate,
    MessagesPlaceholder
)
from langchain.schema import SystemMessage


class PromptFactory:
    @staticmethod
    def create_few_shot_prompt()-> FewShotChatMessagePromptTemplate:
        """prompt for 進階問題

        Returns:
            _type_: FewShotChatMessagePromptTemplate
        """
        examples = [
            {
                "input": "我手上只有企丸丸跟波霸牛，我要如何得到葉泥泥?",
                "output": "1. 企丸丸和波霸牛交配得到草莽豬 2. 再將企丸丸和草莽豬配對，就可以得到葉泥泥。"
            },
        ]

        #ChatPromptTemplate.from_messages 有固定的key , 使用上跟agent所需prompt不符
        example_prompt = ChatPromptTemplate.from_messages(
            [('human', '{input}'), ('ai', '{output}')]
        )

        few_shot_prompt = FewShotChatMessagePromptTemplate(
            examples=examples,
            example_prompt=example_prompt,
        )

        return few_shot_prompt

    @staticmethod
    def create_chat_prompt() -> ChatPromptTemplate:
        
        examples = [
            {
                "input": "我手上只有企丸丸跟波霸牛，我要如何得到葉泥泥?",
                "output": "1. 企丸丸和波霸牛交配得到草莽豬 2. 再將企丸丸和草莽豬配對，就可以得到葉泥泥。"
            },
        ]

        example_prompt = ChatPromptTemplate.from_messages(
            [('human', '{input}'), ('ai', '{output}')]
        )

        few_shot_prompt = FewShotChatMessagePromptTemplate(
            examples=examples,
            example_prompt=example_prompt,
        )
        
                
        prompt = ChatPromptTemplate(
            messages=[
                SystemMessage(content=(
                    "你可以Access db，db裡面存放帕魯的資料以及帕魯的配種資料"
                    "breed table 裡面存放的是一種名叫“帕魯” 生物的配種表，表示parent1跟parent2欄位的帕魯可以生下child欄位的帕魯"
                    "如果使用者問題是跟帕魯相關的，請確實使用run_sqlite_query 'function' 執行SQL查詢語句，並且根據查詢結果作回答"
                    "請先使用'get_example_tool' function 查看範例，裡面有user問題跟要使用甚麼樣的SQL做查詢"
                    "執行任何SQL查詢時，SQL指令請一律加上limit 10，避免過多資料回傳 "
                    "請不要自己生成答案 "
                    "若要查詢帕魯的ID，請去paw_id_name table查詢 "
                    "若要查詢配種資料，請去breed table查詢 "
                    "如果是配種相關問題，可能會需要多個配種步驟來得到答案 "
                    "如果遇到 no such column 錯誤，請先使用describe_tables' function 查詢欄位 "
                    "确保只返回与问题直接相关的数据。在使用工具时 确保只返回与问题直接相关的数据。在使用工具时，必须遵循操作规范，避免执行可能改变数据库状态的操作（如插入、更新、删除等）。"
                    "若问题与数据库内容无关，则直接回答“我不知道” 不要自己生成答案"
                    "回復請用中文。"
                    "下面是一些問題示例以及該問題應該要使用的SQL查詢。"
                )),
                MessagesPlaceholder(variable_name="chat_history"),
                HumanMessagePromptTemplate.from_template("{input}"),
                MessagesPlaceholder(variable_name="agent_scratchpad")
            ]
        )
        return prompt
    
    @staticmethod
    def create_chat_prompt_no_memory() -> ChatPromptTemplate:                
        prompt = ChatPromptTemplate(
            messages=[
                SystemMessage(content=(
                    "你可以Access db，db裡面存放帕魯的資料以及帕魯的配種資料"
                    "breed table 裡面存放的是一種名叫“帕魯” 生物的配種表，表示parent1跟parent2欄位的帕魯可以生下child欄位的帕魯"
                    "請先使用'get_example_tool' function 查看範例，裡面有user問題跟要使用甚麼樣的SQL做查詢 "
                    "請跟據使用者提出的問題來修改Query的參數，不要照抄範例的Query "
                    "再使用run_sqlite_query 'function' 執行SQL查詢語句，並且根據查詢結果作回答 "
                    "執行任何SQL查詢時，SQL指令請一律加上limit 10，避免過多資料回傳 "
                )),
                HumanMessagePromptTemplate.from_template("{input}"),
                MessagesPlaceholder(variable_name="agent_scratchpad")
            ]
        )
        return prompt
    
    # def create_chat_promptV2(memoryId = None, messages = None) -> ChatPromptTemplate:
    #     """Create sys prompt fo paw ai

    #     Args:
    #         messages string : sys prompt string

    #     Returns:
    #         ChatPromptTemplate: sys prompt datastructure for chat model
    #     """
        
    #     messages = [PromptFactory.create_paw_sys_prompt()]
    #     # few_shot_prompt = PromptFactory.create_few_shot_prompt()
    #     # messages.append(few_shot_prompt)
        
    #     if memoryId is not None:
    #         messages.append(MessagesPlaceholder(variable_name=memoryId))
            
    #     messages.append(HumanMessagePromptTemplate.from_template("{input}"))
    #     messages.append(MessagesPlaceholder(variable_name="agent_scratchpad"))
        
        
        
    #     return ChatPromptTemplate(messages)
        
    @staticmethod
    def create_paw_sys_prompt() -> SystemMessage:
        """Create sys prompt fo paw ai

        Returns:
            ChatPromptTemplate: SystemMessage
        """
        return SystemMessage(content=(
                    "你可以Access db，db裡面存放帕魯的資料以及帕魯的配種資料"
                    "breed table 裡面存放的是一種名叫“帕魯” 生物的配種表，表示parent1跟parent2欄位的帕魯可以生下child欄位的帕魯"
                    "請先使用'get_example_tool' function 查看範例，裡面有user問題跟要使用甚麼樣的SQL做查詢"
                    "如果使用者問題是跟帕魯相關的，請確實使用run_sqlite_query 'function' 執行SQL查詢語句，並且根據查詢結果作回答"
                    "執行任何SQL查詢時，SQL指令請一律加上limit 10，避免過多資料回傳 "
                    "請不要自己生成答案 "
                    "若要查詢帕魯的ID，請去paw_id_name table查詢 "
                    "若要查詢配種資料，請去breed table查詢 "
                    "如果是配種相關問題，可能會需要多個配種步驟來得到答案 "
                    "如果遇到 no such column 錯誤，請先使用describe_tables' function 查詢欄位 "
                    "确保只返回与问题直接相关的数据。在使用工具时 确保只返回与问题直接相关的数据。在使用工具时，必须遵循操作规范，避免执行可能改变数据库状态的操作（如插入、更新、删除等）。"
                    "若问题与数据库内容无关，则直接回答“我不知道” 不要自己生成答案"
                    "回復請用中文。"
                    "下面是一些問題示例以及該問題應該要使用的SQL查詢。"
                ))
            
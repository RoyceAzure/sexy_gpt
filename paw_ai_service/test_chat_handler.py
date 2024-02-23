import threading
import time
import pytest
from appserver.server import create_app





@pytest.fixture
def client():
    app = create_app()
    app.config['TESTING'] = True
    with app.test_client() as client:
        yield client





# def test_chain_chat(client):
#     threads = []
#     for _ in range(10):  # 創建10個線程
#         thread = threading.Thread(target=chain_chat)
#         thread.start()
#         threads.append(thread)
    
#     for thread in threads:
#         thread.join()  # 等待所有線程完成

def test_chain_chat(client):
    # app = client.application 
    def chain_chat():
        rv = client.get('/api/v1/chat/init')
        json_data = rv.get_json()
        code = json_data.get("code")
        session_id = json_data.get("sessionid")
        assert code
        assert session_id
        assert rv.status_code == 200

        ques = [
            "魯米兒和冥鎧蠍可以配種出甚麼?",
            "如何生出阿努比斯",
            "有沒有鋼鐵人?",
            "如何生出皮皮雞",
            "有哪些方式可以得到阿努比斯??"
            ]


        for que in ques:
            body = {
                "session_id" : session_id,
                "question" : que
            }
            res = client.post("/api/v1/freechat", json = body)
            json_data = res.get_json()
            code = json_data.get("code")
            assert code
            assert code != 403
        
    threads = []
    for _ in range(10):  # 創建10個線程
        thread = threading.Thread(target=chain_chat)
        thread.daemon = True
        thread.start()
        threads.append(thread)
    time.sleep(200000)
    for thread in threads:
        thread.join()  # 等待所有線程完成
        

    
    
    
def init_session(client):
    
        rv = client.get('/api/v1/chat/init')
        json_data = rv.get_json()
        code = json_data.get("code")
        session_id = json_data.get("sessionid")
        assert code
        assert session_id
        assert rv.status_code == 200

def free_chat(client, session_id):
    
    ques = [
            "魯米兒和冥鎧蠍可以配種出甚麼?",
            "如何生出阿努比斯",
            "有沒有鋼鐵人?",
            "如何生出皮皮雞",
            "有哪些方式可以得到阿努比斯??"
            ]
    
    
    for que in ques:
        body = {
            "session_id" : session_id,
            "question" : que
        }
        res = client.post("/api/v1/freechat", json_data = body)
        code = res.get("code")
        assert code
        assert code == 200

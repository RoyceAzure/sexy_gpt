
# 從自定義的模塊 `common.ma` 中導入 Marshmallow 的實例 `ma`。
# `validate` 和 `ValidationError` 用於自定義驗證邏輯和處理驗證錯誤。
from common.ma import ma
from marshmallow import validate, ValidationError

# 定義一個名為 ReportSchema 的 Schema，這個 Schema 包含了兩個字符串字段：report_id 和 focus。
class ReportSchema(ma.Schema):
    report_id = ma.Str()  # 字符串字段，用於儲存報告的 ID。
    focus = ma.Str(allow_none=True)      # 字符串字段，用於儲存關注的焦點，允許值為 null。
# 定義一個名為 QuesSchema 的 Schema，用於表示某種問題的數據結構。
class QuesSchema(ma.Schema):
    disease = ma.Str()     # 字符串字段，用於儲存疾病的名稱。
    subject_id = ma.List(ma.Str())  # 字符串字段，用於儲存主題的 ID。
    report = ma.List(ma.Nested(ReportSchema))  # 使用 ma.List 來定義一個列表，列表中的每個元素都由 ReportSchema 定義。
# 定義一個名為 QuesSchema 的 Schema，For Freetalk。
class FreetalkSchema(ma.Schema):
    session_id = ma.Str(required=True)
    question = ma.Str()  # 字符串字段，用於儲存主題的 ID。
    # report = ma.List(ma.Nested(ReportSchema))  # 使用 ma.List 來定義一個列表，列表中的每個元素都由 ReportSchema 定義。
# 定義一個名為 StatusSchema 的 Schema，這個 Schema 包含了一個整數字段和一個字符串字段。
class StatusSchema(ma.Schema):
    code = ma.Int()    # 整數字段，用於儲存狀態碼。
    message = ma.Str()  # 字符串字段，用於儲存相關的消息或描述。

# 定義一個名為 ResultSchema 的 Schema，這個 Schema 包含了一個整數字段和一個字符串字段。
class ResultSchema(ma.Schema):
    report_id = ma.Int()  # 整數字段，用於儲存報告的 ID。
    focus = ma.Str()      # 字符串字段，用於儲存關注的焦點。

# 定義一個名為 AnsSchema 的 Schema，用於表示某種答案的數據結構。
class AnsSchema(ma.Schema):
    # status = ma.Nested(StatusSchema)  # 嵌套字段，使用 StatusSchema 來序列化/反序列化 `status` 字段。
    # result = ma.Nested(ResultSchema)  # 嵌套字段，使用 ResultSchema 來序列化/反序列化 `result` 字段。
    ans = ma.Str()                    
    # token = ma.Int()                  # 整數字段，用於儲存某種標記或識別符。
    # referencetext = ma.Str()

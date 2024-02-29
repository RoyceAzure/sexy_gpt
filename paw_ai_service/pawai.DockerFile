# 使用官方 Python 基礎映像
FROM python:3.11-slim

# 設置工作目錄
WORKDIR /app

# 將當前目錄下的所有文件複製到容器的 /app 目錄
COPY . /app

# 安裝 requirements.txt 中列出的所有依賴
RUN pip install --no-cache-dir -r requirements.txt

# 暴露 Flask 預設端口 5000
EXPOSE 5000
ENV FLASK_APP=app.py
ENV FLASK_RUN_HOST=0.0.0.0

CMD ["python", "app.py"]

FROM pytorch/pytorch:2.5.1-cuda12.1-cudnn9-devel

WORKDIR /app
ENV PYTHONUNBUFFERED=1

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt \
	-i https://pypi.tuna.tsinghua.edu.cn/simple \
	--timeout 120

COPY . .

RUN mkdir -p /app/log /app/save_model
VOLUME ["/app/log", "/app/save_model"]

CMD ["python", "main.py"]

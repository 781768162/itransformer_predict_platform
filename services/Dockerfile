FROM pytorch/pytorch:2.5.1-cpu

WORKDIR /app
ENV PYTHONUNBUFFERED=1

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

RUN mkdir -p /app/log /app/save_model
VOLUME ["/app/log", "/app/save_model"]

CMD ["python", "main.py"]

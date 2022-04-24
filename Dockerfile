FROM python:alpine3.15

WORKDIR /app

COPY requirements.txt .
COPY main.py .

RUN apk update
RUN pip install --no-cache-dir --upgrade -r requirements.txt

RUN rm -rf /app/requirements.txt

ENTRYPOINT ["uvicorn", "main:APP", "--host", "0.0.0.0", "--port", "80"]
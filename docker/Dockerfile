# FROM python:3.9-slim

WORKDIR /app

COPY ./docker/requirements.txt .

RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir -r requirements.txt

COPY ./src .

EXPOSE 5500

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "5500"]
FROM python:3.8-slim

WORKDIR /app

RUN apt-get update && apt-get install -y g++
COPY /src/recommendation/requirements.txt /app/requirements.txt
RUN pip install -r requirements.txt

COPY /src/recommendation /app

CMD ["python", "-m", "flask", "run", "--host=0.0.0.0", "--port=80"]

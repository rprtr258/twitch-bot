FROM python:3.9.14-buster

WORKDIR /app
COPY requirements.txt .
RUN pip3 install -r requirements.txt

COPY *.py /app/
ENTRYPOINT [ "timeout", "5", "python3", "pyth.py", "-c" ]

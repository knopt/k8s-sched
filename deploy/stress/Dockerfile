FROM python:3.7

RUN mkdir -p /python/workdir
WORKDIR /python/workdir
RUN pip3 install --upgrade pip setuptools && pip3 install scipy
RUN apt-get update -y && apt-get install stress-ng -y
RUN echo $(which stress-ng)

COPY stresscpu.py stresscpu.py

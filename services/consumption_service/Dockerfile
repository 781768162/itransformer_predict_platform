FROM nvidia/cuda:12.1.1-cudnn8-runtime-ubuntu22.04 AS stage0

SHELL ["/bin/bash", "-l", "-c"]

RUN apt-get -o Acquire::http::Proxy="http://mirrors.aliyun.com/ubuntu/" update && \
    apt-get -o Acquire::http::Proxy="http://mirrors.aliyun.com/ubuntu/" install -y \
        python3.10 \
        python3-pip \
    apt-get clean

RUN ln -sf /usr/bin/python3 /usr/bin/python

FROM stage0 AS stage1


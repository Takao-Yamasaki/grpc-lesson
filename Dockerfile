FROM golang:1.20.5

RUN apt-get update && apt-get install -y \
    less \
    protobuf-compiler \
    tree \
    && apt-get clean \
    && rm -rf /var/lib/apt/list/*

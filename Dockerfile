FROM golang:1.20.5

RUN apt-get update && apt-get install -y \
    less \
    protobuf-compiler \
    tree \
    # ダウンロードしたパッケージファイルを削除
    && apt-get clean \
    # aptのパッケージキャッシュを削除
    && rm -rf /var/lib/apt/list/*

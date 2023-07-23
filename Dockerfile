FROM golang:1.20.5

RUN apt-get update && apt-get install -y \
    less \
    protobuf-compiler \
    tree \
    make \
    # ダウンロードしたパッケージファイルを削除
    && apt-get clean \
    # aptのパッケージキャッシュを削除
    && rm -rf /var/lib/apt/list/*

WORKDIR /go/src/workspace

COPY Makefile go.mod go.sum /go/src/workspace/

# Protobuf用のGoのパッケージをインストール
RUN make install

# モジュールの依存関係をダウンロード
RUN go mod download

CMD [ "/bin/bash" ]



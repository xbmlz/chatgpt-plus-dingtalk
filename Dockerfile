FROM golang:1.18.10-alpine3.16 AS builder

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o chatgpt-plus-dingtalk .

FROM alpine:3.16

ARG TZ="Asia/Shanghai"

ENV TZ ${TZ}

RUN mkdir /app && apk upgrade \
    && apk add bash tzdata \
    && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone

WORKDIR /app
COPY --from=builder /app/ .
RUN chmod +x chatgpt-plus-dingtalk && cp config.example.yaml config.yaml

CMD ./chatgpt-plus-dingtalk
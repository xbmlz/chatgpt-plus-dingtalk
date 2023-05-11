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
    && echo ${TZ} > /etc/timezone \
    && apk add --no-cache chromium \
    && apk add --update ttf-dejavu fontconfig \
    && rm -rf /var/cache/apk/* \
    && mkfontscale && mkfontdir && fc-cache -fv

WORKDIR /app
COPY --from=builder /app/ .
COPY pkg/mermaid/SIMSUN.TTC /usr/share/fonts/zh/
RUN chmod +x chatgpt-plus-dingtalk && cp config.example.yml config.yml

CMD ./chatgpt-plus-dingtalk
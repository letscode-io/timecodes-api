FROM golang:alpine AS builder

ENV APP_HOME=/usr/src/app

RUN apk update \
    && apk upgrade \
    && apk add --upgrade git \
    && rm -rf /var/cache/apk/*

RUN mkdir -p $APP_HOME

WORKDIR $APP_HOME

COPY go.mod go.sum $APP_HOME/

RUN go mod download

COPY . $APP_HOME

WORKDIR $APP_HOME/cmd
RUN go build -o $APP_HOME/application

# Final stage

FROM alpine

ENV APP_HOME=/usr/src/app

RUN apk update --purge \
    && apk upgrade --purge \
    && rm -rf /var/cache/apk/*

WORKDIR $APP_HOME

COPY --from=builder $APP_HOME/application $APP_HOME/

ENTRYPOINT ["./application"]

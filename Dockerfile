FROM docker.pkg.github.com/letscode-io/timecodes-api/timecodes-api_base:latest AS base

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

RUN CGO_ENABLED=0 GOOS=linux go build -o $APP_HOME/application $APP_HOME/cmd

# Final stage

FROM alpine AS final

ENV APP_HOME=/usr/src/app

RUN apk update --purge \
    && apk upgrade --purge \
    && rm -rf /var/cache/apk/*

WORKDIR $APP_HOME

COPY --from=base $APP_HOME/application $APP_HOME/

ENTRYPOINT ["./application"]

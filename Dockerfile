FROM golang:1.21.3
WORKDIR /atri
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o ./out .

FROM alpine:3.18
WORKDIR /atri
COPY --from=0 /atri/out /bin/atri

ENV ENV=prod
ARG TOKEN_DISCORD_APPLICATION
ENV TOKEN_DISCORD_APPLICATION=${TOKEN_DISCORD_APPLICATION}
CMD ["atri"]

FROM golang:1.21.5
WORKDIR /atri
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o ./out .

FROM alpine:3.19
WORKDIR /atri
COPY --from=0 /atri/out /bin/atri
CMD ["atri"]

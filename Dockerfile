FROM golang:alpine as builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN go build -p 2 -o app cmd/kafka-test/main.go cmd/kafka-test/flags.go

FROM alpine 

WORKDIR /build

COPY --from=builder /build/app /build/app
COPY  /migrations /build/migrations
RUN touch /build/logs.json

CMD ["./app"]
FROM golang:1.20.7-alpine3.18 as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/app ./cmd/main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /usr/local/bin/app /app/cyclotron-user-service

ENTRYPOINT ["/app/cyclotron-user-service"]



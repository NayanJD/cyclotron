FROM alpine:3.18

ARG GOOSARCH=linux-amd64

WORKDIR /app

COPY ./dist/userSvc-$GOOSARCH /app/userSvc

ENTRYPOINT ["/app/userSvc"]



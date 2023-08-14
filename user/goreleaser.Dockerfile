ARG BINARY_NAME=cyclotron-user-service

FROM alpine:3.18

WORKDIR /app

COPY ./$BINARY_NAME /app/BINARY_NAME

ENTRYPOINT ["/app/${BINARY_NAME}"]



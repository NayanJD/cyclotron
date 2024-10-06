## Cyclotron: Collection of Go Microservices built using [go-kit](https://github.com/go-kit/kit)

What's this project is for:

1. Go microservices reference following all production standard practises
2. Benchmark/load test go microservices

### Building

This project uses [Taskfile](https://taskfile.dev/) to generate all sorts for artifacts.

To build just the binary of any service (taking example of user service):

```
GOOS=linux GOOS=amd64 task build-usersvc
```

To build local docker image of a service:

```
GOOS=linux GOARCH=arm64 task build-usersvc-docker-local
```

To build k6 binary with dependecies:

```
GOOS=linux GOARCH=arm64 task build-k6
```

### Code Generation

This project uses [kit CLI](https://github.com/GrantZheng/kit) to generate services.


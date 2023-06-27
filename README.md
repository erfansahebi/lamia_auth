# Lamia

Project with microservice architecture.

## Project Overview

Microservices with API Gateway which handles incoming HTTP requests.
HTTP requests will be forwarded to these Microservices by <b>gRPC</b>.
Additionally, we deal with Authentication, Redis and PostgreSQL.

## Architecture

Microservice architecture is an architectural style that structures an application as a collection of small, loosely
coupled, and independently deployable services.
In this architecture, each service represents a specific business capability and can be developed, deployed, and scaled
independently.
Services communicate with each other through well-defined APIs, typically over lightweight protocols like HTTP or
messaging queues.

gRPC is a high-performance, open-source framework developed by Google for building remote procedure call (RPC) systems.
It uses Protocol Buffers (protobuf) as the interface definition language (IDL) for describing the API contract between the client and server.
gRPC supports multiple programming languages, including Go, Python, Java, and more.

## Prerequisites

- Git
- Docker

## Getting Started

Repositories:

- [Lamia EnvDev](https://github.com/erfansahebi/lamia_envdev)
- [Lamia Gateway](https://github.com/erfansahebi/lamia_gateway)
- [Lamia Auth](https://github.com/erfansahebi/lamia_auth)
- [Lamia Shared](https://github.com/erfansahebi/lamia_shared)
- [Lamia Nginx](https://github.com/erfansahebi/lamia_nginx)

<hr>

### 1. Install git
 ```shell
sudo apt install git-all
```

### 2. Pull [Lamia EnvDev](https://github.com/erfansahebi/lamia_envdev)
```shell
git pull https://github.com/erfansahebi/lamia_envdev
```
In this repo, we have a shell for pulling all Repositories.

### 3. Pull all Repositories:
```shell
chmod +x build.sh
```
```shell
sh build.sh
```

### 4. Make Docker Images
```shell
make build
```

### 5. Set Environments
```shell
cp .env.example .env
```

### 6. run Docker:
```shell
dokcer compose up --build -d
```
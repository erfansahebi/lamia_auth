# syntax=docker/dockerfile:1.3

FROM golang:1.20-alpine

WORKDIR /src/server

RUN wget -qO- https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

RUN apk add git

# Copy golang dependency manifests
COPY go.mod .
COPY go.sum .

# Cache the downloaded dependency in the layer.
RUN --mount=type=cache,target=/root/.cache/go-build go mod download

COPY . /src/server
RUN --mount=type=cache,target=/root/.cache/go-build \
go get -d -v ./ && \
go install -v ./

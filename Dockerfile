FROM golang:latest AS GO_BUILD
COPY . .
ENV GOPATH=/
CMD go run ./cmd/main/main.go
USER root
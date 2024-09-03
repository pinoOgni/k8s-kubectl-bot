FROM golang:1.22.1 AS builder

RUN apt-get update && apt-get install -y curl && \
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go main.go

# Ensure correct architecture
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build -o k8s-kubectl-bot

FROM alpine:latest

# Install ca-certificates and other dependencies for kubectl
RUN apk add --no-cache ca-certificates

COPY --from=builder /usr/local/bin/kubectl /usr/local/bin/kubectl

COPY --from=builder /app/k8s-kubectl-bot /k8s-kubectl-bot

ENTRYPOINT ["/k8s-kubectl-bot"]

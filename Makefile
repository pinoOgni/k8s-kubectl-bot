# Define a default version
VERSION ?= 0.1.0

.PHONY: all image build

all: build

image-build:
	docker build -t pinoogni/k8s-kubectl-bot:$(VERSION) .

image-push:
	docker push pinoogni/k8s-kubectl-bot:$(VERSION)

build:
	go build -o bin/k8s-kubectl-bot .

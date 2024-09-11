
REGISTRY ?= quay.io/dvossel
TAG ?= latest
IMAGE_NAME ?= kubevirt-cpu-model-cluster-analyzer

.PHONY: build
build:
	go build -o kubevirt-cpu-model-cluster-analyzer main.go

.PHONY: docker-build
docker-build:
	docker build . -t $(IMAGE_NAME)):$(TAG) --file Dockerfile

.PHONY: docker-push
docker-push:
	docker push $(IMAGE_NAME):$(TAG)

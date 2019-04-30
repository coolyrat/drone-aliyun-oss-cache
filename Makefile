DOCKER_IMAGE := luischan/drone-oss-cache:latest
GIT_COMMIT := $(shell git rev-parse --short HEAD)
LDFLAGS := -X oa.eekui.group/oa-suite/pkg/version.GitVersion=$(GIT_COMMIT) \
	-X 'oa.eekui.group/oa-suite/pkg/version.GoVersion=$(shell go version)'


.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o drone-oss-cache

.PHONY: docker
docker: build
	docker build -t $(DOCKER_IMAGE) .

.PHONY: push
push: docker
	docker push $(DOCKER_IMAGE)
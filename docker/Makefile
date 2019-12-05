#Makefile related to docker

DOCKER_IMAGE ?= namechain
DOCKER_TAG ?= 0.1
DOCKER_REF := $(DOCKER_IMAGE):$(DOCKER_TAG)


.PHONY: docker
docker:
	docker build -t $(DOCKER_REF) .

.PHONY: docker-build
docker-build:
	docker run -ti --rm -v $(CURDIR):/srv/app/src/code.gitea.io/gitea -w /srv/app/src/code.gitea.io/gitea -e TAGS="bindata $(TAGS)" webhippie/golang:edge make clean generate build

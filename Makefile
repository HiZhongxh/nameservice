PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=NameService \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=nsd \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=nscli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 

BUILD_FLAGS := -ldflags '$(ldflags)'

include Makefile.ledger
all: build

build: go.sum
	go build -mod=readonly $(BUILD_FLAGS) -o build/nsd ./cmd/nsd
	go build -mod=readonly $(BUILD_FLAGS) -o build/nscli ./cmd/nscli

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/nsd
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/nscli

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

proto:
	@echo "\033[;32mbuild protobuf file \033[0m"
	make -C x/nameservice/internal/types/pb

build-linux: go.sum
    LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

build-docker-nsdnode:
	$(MAKE) -C networks/local

# Stop testnet
localnet-stop:
	docker-compose down

# Run a 4-node testnet locally
localnet-start: build-linux localnet-stop
	@if ! [ -f build/node0/nsd/config/genesis.json ]; then docker run -d -v $(CURDIR)/build:/nsd:Z tendermint/nsdnode -o 127.0.0.1; fi
	docker-compose up -d

test:
	@go test -mod=readonly $(PACKAGES)
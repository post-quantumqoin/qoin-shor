SHELL=/usr/bin/env bash

all: build
.PHONY: all

unexport GOFLAGS

GOCC?=go

GOVERSION:=$(shell $(GOCC) version | tr ' ' '\n' | grep go1 | sed 's/^go//' | awk -F. '{printf "%d%03d%03d", $$1, $$2, $$3}')
GOVERSIONMIN:=$(shell cat GO_VERSION_MIN | awk -F. '{printf "%d%03d%03d", $$1, $$2, $$3}')

ifeq ($(shell expr $(GOVERSION) \< $(GOVERSIONMIN)), 1)
$(warning Your Golang version is go$(shell expr $(GOVERSION) / 1000000).$(shell expr $(GOVERSION) % 1000000 / 1000).$(shell expr $(GOVERSION) % 1000))
$(error Update Golang to version to at least $(shell cat GO_VERSION_MIN))
endif

# git modules that need to be loaded
MODULES:=

CLEAN:=
BINS:=

ldflags=-X=github.com/filecoin-project/lotus/build.CurrentCommit=+git.$(subst -,.,$(shell git describe --always --match=NeVeRmAtCh --dirty 2>/dev/null || git rev-parse --short HEAD 2>/dev/null))
ifneq ($(strip $(LDFLAGS)),)
	ldflags+=-extldflags=$(LDFLAGS)
endif

GOFLAGS+=-ldflags="$(ldflags)"


## FFI

FFI_PATH:=ext/qvm/
FFI_DEPS:=.install-filcrypto
FFI_DEPS:=$(addprefix $(FFI_PATH),$(FFI_DEPS))

$(FFI_DEPS): build/.filecoin-install;

build/.filecoin-install: $(FFI_PATH)
	$(MAKE) -C $(FFI_PATH) $(FFI_DEPS:$(FFI_PATH)%=%)
	@touch $@

MODULES+=$(FFI_PATH)
BUILD_DEPS+=build/.filecoin-install
CLEAN+=build/.filecoin-install

ffi-version-check:
	@[[ "$$(awk '/const Version/{print $$5}' ext/qvm/version.go)" -eq 3 ]] || (echo "FFI version mismatch, update submodules"; exit 1)
BUILD_DEPS+=ffi-version-check

.PHONY: ffi-version-check

$(MODULES): build/.update-modules ;
# dummy file that marks the last time modules were updated
build/.update-modules:
	git submodule update --init --recursive
	touch $@

# end git modules

## MAIN BINARIES

CLEAN+=build/.update-modules

deps: $(BUILD_DEPS)
.PHONY: deps

build-devnets: build shor-seed shor-shed shor-provider
.PHONY: build-devnets

debug: GOFLAGS+=-tags=debug
debug: build-devnets

2k: GOFLAGS+=-tags=2k
2k: build-devnets

calibnet: GOFLAGS+=-tags=calibnet
calibnet: build-devnets

butterflynet: GOFLAGS+=-tags=butterflynet
butterflynet: build-devnets

interopnet: GOFLAGS+=-tags=interopnet
interopnet: build-devnets

shor: $(BUILD_DEPS)
	rm -f shor
	$(GOCC) build $(GOFLAGS) -o shor ./cmd/shor
.PHONY: shor
BINS+=shor

shor-miner: $(BUILD_DEPS)
	rm -f shor-miner
	$(GOCC) build $(GOFLAGS) -o shor-miner ./cmd/shor-miner
.PHONY: shor-miner
BINS+=shor-miner

shor-provider: $(BUILD_DEPS)
	rm -f shor-provider
	$(GOCC) build $(GOFLAGS) -o shor-provider ./cmd/shor-provider
.PHONY: shor-provider
BINS+=shor-provider

lp2k: GOFLAGS+=-tags=2k
lp2k: lotus-provider

shor-worker: $(BUILD_DEPS)
	rm -f shor-worker
	$(GOCC) build $(GOFLAGS) -o shor-worker ./cmd/shor-worker
.PHONY: shor-worker
BINS+=shor-worker

shor-shed: $(BUILD_DEPS)
	rm -f shor-shed
	$(GOCC) build $(GOFLAGS) -o shor-shed ./cmd/shor-shed
.PHONY: shor-shed
BINS+=shor-shed

shor-gateway: $(BUILD_DEPS)
	rm -f shor-gateway
	$(GOCC) build $(GOFLAGS) -o shor-gateway ./cmd/shor-gateway
.PHONY: shor-gateway
BINS+=shor-gateway

build: shor shor-miner shor-worker shor-provider
	@[[ $$(type -P "shor") ]] && echo "Caution: you have \
an existing shor binary in your PATH. This may cause problems if you don't run 'sudo make install'" || true

.PHONY: build

install: shor-daemon shor-miner shor-worker shor-provider

install-daemon:
	install -C ./shor /usr/local/bin/shor

install-miner:
	install -C ./shor-miner /usr/local/bin/shor-miner

install-provider:
	install -C ./shor-provider /usr/local/bin/shor-provider

install-worker:
	install -C ./shor-worker /usr/local/bin/shor-worker

install-app:
	install -C ./$(APP) /usr/local/bin/$(APP)

uninstall: uninstall-daemon uninstall-miner uninstall-worker
.PHONY: uninstall

uninstall-daemon:
	rm -f /usr/local/bin/shor

uninstall-miner:
	rm -f /usr/local/bin/shor-miner

uninstall-provider:
	rm -f /usr/local/bin/shor-provider

uninstall-worker:
	rm -f /usr/local/bin/shor-worker

# TOOLS

shor-seed: $(BUILD_DEPS)
	rm -f shor-seed
	$(GOCC) build $(GOFLAGS) -o shor-seed ./cmd/shor-seed

.PHONY: shor-seed
BINS+=shor-seed

benchmarks:
	$(GOCC) run github.com/whyrusleeping/bencher ./... > bench.json
	@echo Submitting results
	@curl -X POST 'http://benchmark.kittyhawk.wtf/benchmark' -d '@bench.json' -u "${benchmark_http_cred}"
.PHONY: benchmarks

shor-fountain:
	rm -f shor-fountain
	$(GOCC) build $(GOFLAGS) -o shor-fountain ./cmd/shor-fountain
	$(GOCC) run github.com/GeertJohan/go.rice/rice append --exec shor-fountain -i ./cmd/lotus-fountain -i ./build
.PHONY: shor-fountain
BINS+=shor-fountain

shor-bench:
	rm -f shor-bench
	$(GOCC) build $(GOFLAGS) -o shor-bench ./cmd/shor-bench
.PHONY: shor-bench
BINS+=shor-bench

shor-stats:
	rm -f shor-stats
	$(GOCC) build $(GOFLAGS) -o shor-stats ./cmd/shor-stats
.PHONY: shor-stats
BINS+=shor-stats

shor-pcr:
	rm -f shor-pcr
	$(GOCC) build $(GOFLAGS) -o shor-pcr ./cmd/shor-pcr
.PHONY: shor-pcr
BINS+=shor-pcr

shor-health:
	rm -f shor-health
	$(GOCC) build -o shor-health ./cmd/shor-health
.PHONY: shor-health
BINS+=shor-health

shor-wallet: $(BUILD_DEPS)
	rm -f shor-wallet
	$(GOCC) build $(GOFLAGS) -o shor-wallet ./cmd/shor-wallet
.PHONY: shor-wallet
BINS+=shor-wallet

shor-keygen:
	rm -f shor-keygen
	$(GOCC) build -o shor-keygen ./cmd/shor-keygen
.PHONY: shor-keygen
BINS+=shor-keygen

testground:
	$(GOCC) build -tags testground -o /dev/null ./cmd/shor
.PHONY: testground
BINS+=testground


tvx:
	rm -f tvx
	$(GOCC) build -o tvx ./cmd/tvx
.PHONY: tvx
BINS+=tvx

shor-sim: $(BUILD_DEPS)
	rm -f shor-sim
	$(GOCC) build $(GOFLAGS) -o shor-sim ./cmd/shor-sim
.PHONY: shor-sim
BINS+=shor-sim

# SYSTEMD

install-daemon-service: install-daemon
	mkdir -p /etc/systemd/system
	mkdir -p /var/log/lotus
	install -C -m 0644 ./scripts/lotus-daemon.service /etc/systemd/system/lotus-daemon.service
	systemctl daemon-reload
	@echo
	@echo "lotus-daemon service installed. Don't forget to run 'sudo systemctl start lotus-daemon' to start it and 'sudo systemctl enable lotus-daemon' for it to be enabled on startup."

install-miner-service: install-miner install-daemon-service
	mkdir -p /etc/systemd/system
	mkdir -p /var/log/lotus
	install -C -m 0644 ./scripts/lotus-miner.service /etc/systemd/system/lotus-miner.service
	systemctl daemon-reload
	@echo
	@echo "lotus-miner service installed. Don't forget to run 'sudo systemctl start lotus-miner' to start it and 'sudo systemctl enable lotus-miner' for it to be enabled on startup."

install-provider-service: install-provider install-daemon-service
	mkdir -p /etc/systemd/system
	mkdir -p /var/log/lotus
	install -C -m 0644 ./scripts/lotus-provider.service /etc/systemd/system/lotus-provider.service
	systemctl daemon-reload
	@echo
	@echo "lotus-provider service installed. Don't forget to run 'sudo systemctl start lotus-provider' to start it and 'sudo systemctl enable lotus-provider' for it to be enabled on startup."

install-main-services: install-miner-service

install-all-services: install-main-services

install-services: install-main-services

clean-daemon-service: clean-miner-service
	-systemctl stop lotus-daemon
	-systemctl disable lotus-daemon
	rm -f /etc/systemd/system/lotus-daemon.service
	systemctl daemon-reload

clean-miner-service:
	-systemctl stop lotus-miner
	-systemctl disable lotus-miner
	rm -f /etc/systemd/system/lotus-miner.service
	systemctl daemon-reload

clean-provider-service:
	-systemctl stop lotus-provider
	-systemctl disable lotus-provider
	rm -f /etc/systemd/system/lotus-provider.service
	systemctl daemon-reload

clean-main-services: clean-daemon-service

clean-all-services: clean-main-services

clean-services: clean-all-services

# MISC

buildall: $(BINS)

install-completions:
	mkdir -p /usr/share/bash-completion/completions /usr/local/share/zsh/site-functions/
	install -C ./scripts/bash-completion/lotus /usr/share/bash-completion/completions/lotus
	install -C ./scripts/zsh-completion/lotus /usr/local/share/zsh/site-functions/_lotus

clean:
	rm -rf $(CLEAN) $(BINS)
	-$(MAKE) -C $(FFI_PATH) clean
.PHONY: clean

dist-clean:
	git clean -xdff
	git submodule deinit --all -f
.PHONY: dist-clean

type-gen: api-gen
	$(GOCC) run ./gen/main.go
	$(GOCC) generate -x ./...
	goimports -w api/

actors-code-gen:
	$(GOCC) run ./gen/inline-gen . gen/inlinegen-data.json
	$(GOCC) run ./chain/actors/agen
	$(GOCC) fmt ./...

actors-gen: actors-code-gen 
	./scripts/fiximports
.PHONY: actors-gen

bundle-gen:
	$(GOCC) run ./gen/bundle $(VERSION) $(RELEASE) $(RELEASE_OVERRIDES)
	$(GOCC) fmt ./build/...
.PHONY: bundle-gen


api-gen:
	$(GOCC) run ./gen/api
	goimports -w api
	goimports -w api
.PHONY: api-gen

cfgdoc-gen:
	$(GOCC) run ./node/config/cfgdocgen > ./node/config/doc_gen.go

appimage: lotus
	rm -rf appimage-builder-cache || true
	rm AppDir/io.filecoin.lotus.desktop || true
	rm AppDir/icon.svg || true
	rm Appdir/AppRun || true
	mkdir -p AppDir/usr/bin
	cp ./lotus AppDir/usr/bin/
	appimage-builder

docsgen: docsgen-md docsgen-openrpc fiximports

docsgen-md-bin: api-gen actors-gen
	$(GOCC) build $(GOFLAGS) -o docgen-md ./api/docgen/cmd
docsgen-openrpc-bin: api-gen actors-gen
	$(GOCC) build $(GOFLAGS) -o docgen-openrpc ./api/docgen-openrpc/cmd

docsgen-md: docsgen-md-full docsgen-md-storage docsgen-md-worker docsgen-md-provider

docsgen-md-full: docsgen-md-bin
	./docgen-md "api/api_full.go" "FullNode" "api" "./api" > documentation/en/api-v1-unstable-methods.md
	./docgen-md "api/v0api/full.go" "FullNode" "v0api" "./api/v0api" > documentation/en/api-v0-methods.md
docsgen-md-storage: docsgen-md-bin
	./docgen-md "api/api_storage.go" "StorageMiner" "api" "./api" > documentation/en/api-v0-methods-miner.md
docsgen-md-worker: docsgen-md-bin
	./docgen-md "api/api_worker.go" "Worker" "api" "./api" > documentation/en/api-v0-methods-worker.md
docsgen-md-provider: docsgen-md-bin
	./docgen-md "api/api_lp.go" "Provider" "api" "./api" > documentation/en/api-v0-methods-provider.md

docsgen-openrpc: docsgen-openrpc-full docsgen-openrpc-storage docsgen-openrpc-worker docsgen-openrpc-gateway

docsgen-openrpc-full: docsgen-openrpc-bin
	./docgen-openrpc "api/api_full.go" "FullNode" "api" "./api" -gzip > build/openrpc/full.json.gz
docsgen-openrpc-storage: docsgen-openrpc-bin
	./docgen-openrpc "api/api_storage.go" "StorageMiner" "api" "./api" -gzip > build/openrpc/miner.json.gz
docsgen-openrpc-worker: docsgen-openrpc-bin
	./docgen-openrpc "api/api_worker.go" "Worker" "api" "./api" -gzip > build/openrpc/worker.json.gz
docsgen-openrpc-gateway: docsgen-openrpc-bin
	./docgen-openrpc "api/api_gateway.go" "Gateway" "api" "./api" -gzip > build/openrpc/gateway.json.gz

.PHONY: docsgen docsgen-md-bin docsgen-openrpc-bin

fiximports:
	./scripts/fiximports

gen: actors-code-gen type-gen cfgdoc-gen docsgen api-gen circleci
	./scripts/fiximports
	@echo ">>> IF YOU'VE MODIFIED THE CLI OR CONFIG, REMEMBER TO ALSO RUN 'make docsgen-cli'"
.PHONY: gen

jen: gen

snap: lotus lotus-miner lotus-worker lotus-provider
	snapcraft
	# snapcraft upload ./lotus_*.snap

# separate from gen because it needs binaries
docsgen-cli: lotus lotus-miner lotus-worker lotus-provider
	python3 ./scripts/generate-lotus-cli.py
	./lotus config default > documentation/en/default-lotus-config.toml
	./lotus-miner config default > documentation/en/default-lotus-miner-config.toml
	./lotus-provider config default > documentation/en/default-lotus-provider-config.toml
.PHONY: docsgen-cli

print-%:
	@echo $*=$($*)

circleci:
	go generate -x ./.circleci

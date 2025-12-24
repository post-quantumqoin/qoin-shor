module github.com/post-quantumqoin/qoin-shor

go 1.24.6

toolchain go1.24.11

require (
	contrib.go.opencensus.io/exporter/prometheus v0.4.2
	github.com/dgraph-io/badger/v2 v2.2007.4
	github.com/docker/go-units v0.5.0
	github.com/dustin/go-humanize v1.0.1
	github.com/elastic/gosigar v0.14.2
	github.com/filecoin-project/kubo-api-client v0.27.0
	github.com/golang/mock v1.6.0
	github.com/gorilla/websocket v1.5.1
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/influxdata/influxdb1-client v0.0.0-00010101000000-000000000000
	github.com/ipfs/boxo v0.32.0
	github.com/ipfs/go-block-format v0.2.2
	github.com/ipfs/go-cid v0.6.0
	github.com/ipfs/go-datastore v0.6.0
	github.com/ipfs/go-ipld-format v0.6.2
	github.com/ipfs/go-log/v2 v2.5.1
	github.com/ipsn/go-secp256k1 v0.0.0-20180726113642-9d62b9f0bc52
	github.com/libp2p/go-buffer-pool v0.1.0
	github.com/libp2p/go-msgio v0.3.0
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	github.com/multiformats/go-base32 v0.1.0
	github.com/multiformats/go-multiaddr v0.12.2
	github.com/multiformats/go-multihash v0.2.3
	github.com/post-quantumqoin/address v0.0.2
	github.com/post-quantumqoin/core-types v0.3.0
	github.com/post-quantumqoin/go-jsonrpc v0.0.1
	github.com/post-quantumqoin/specs-contracts v0.0.2
	github.com/prometheus/client_golang v1.18.0
	github.com/raulk/clock v1.1.0
	github.com/stretchr/testify v1.11.1
	github.com/whyrusleeping/cbor-gen v0.3.1
	github.com/xorcare/golden v0.8.3
	go.opencensus.io v0.24.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.27.0
	golang.org/x/sync v0.6.0
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/crackcomm/go-gitignore v0.0.0-20231225121904-e25f5bc08668 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/dgraph-io/ristretto v0.0.3-0.20200630154024-f66de99634de // indirect
	github.com/dgryski/go-farm v0.0.0-20190423205320-6a90982ecee2 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/gopacket v1.1.19 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/ipfs/bbloom v0.0.4 // indirect
	github.com/ipfs/go-ipfs-cmds v0.10.0 // indirect
	github.com/ipfs/go-ipld-cbor v0.2.1 // indirect
	github.com/ipfs/go-ipld-legacy v0.2.1 // indirect
	github.com/ipfs/go-log v1.0.5 // indirect
	github.com/ipfs/go-metrics-interface v0.0.1 // indirect
	github.com/ipld/go-codec-dagpb v1.6.0 // indirect
	github.com/ipld/go-ipld-prime v0.21.0 // indirect
	github.com/jbenet/goprocess v0.1.4 // indirect
	github.com/klauspost/compress v1.17.6 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/libp2p/go-cidranger v1.1.0 // indirect
	github.com/libp2p/go-libp2p v0.33.0 // indirect
	github.com/libp2p/go-libp2p-asn-util v0.4.1 // indirect
	github.com/libp2p/go-libp2p-kad-dht v0.24.4 // indirect
	github.com/libp2p/go-libp2p-kbucket v0.6.3 // indirect
	github.com/libp2p/go-libp2p-record v0.2.0 // indirect
	github.com/libp2p/go-libp2p-routing-helpers v0.7.3 // indirect
	github.com/libp2p/go-netroute v0.2.1 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/miekg/dns v1.1.58 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base36 v0.2.0 // indirect
	github.com/multiformats/go-multiaddr-dns v0.3.1 // indirect
	github.com/multiformats/go-multibase v0.2.0 // indirect
	github.com/multiformats/go-multicodec v0.9.0 // indirect
	github.com/multiformats/go-multistream v0.5.0 // indirect
	github.com/multiformats/go-varint v0.1.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/polydawn/refmt v0.89.0 // indirect
	github.com/post-quantumqoin/bitset v0.1.0 // indirect
	github.com/prometheus/client_model v0.6.0 // indirect
	github.com/prometheus/common v0.47.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/prometheus/statsd_exporter v0.22.7 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/samber/lo v1.39.0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/whyrusleeping/base32 v0.0.0-20170828182744-c30ac30633cc // indirect
	github.com/whyrusleeping/go-keyspace v0.0.0-20160322163242-5b898ac5add1 // indirect
	go.opentelemetry.io/otel v1.21.0 // indirect
	go.opentelemetry.io/otel/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.21.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/exp v0.0.0-20240213143201-ec583247a57a // indirect
	golang.org/x/mod v0.15.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/tools v0.18.0 // indirect
	gonum.org/v1/gonum v0.14.0 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	lukechampine.com/blake3 v1.4.1 // indirect
)

replace github.com/elastic/gosigar => github.com/elastic/gosigar v0.14.2

replace github.com/gorilla/websocket => github.com/gorilla/websocket v1.5.1

replace github.com/influxdata/influxdb1-client => github.com/influxdata/influxdb1-client v0.0.0-20200827194710-b269163b24ab

replace github.com/ipfs/boxo => github.com/ipfs/boxo v0.18.0

replace github.com/ipfs/go-cid => github.com/ipfs/go-cid v0.4.1

replace github.com/ipfs/go-datastore => github.com/ipfs/go-datastore v0.6.0

replace github.com/ipfs/go-ipld-format => github.com/ipfs/go-ipld-format v0.6.0

replace github.com/ipfs/go-log/v2 => github.com/ipfs/go-log/v2 v2.5.1

replace github.com/multiformats/go-multiaddr => github.com/multiformats/go-multiaddr v0.12.2

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v1.18.0

replace github.com/stretchr/testify => github.com/stretchr/testify v1.9.0

replace github.com/whyrusleeping/cbor-gen => github.com/whyrusleeping/cbor-gen v0.1.0

replace go.uber.org/zap => go.uber.org/zap v1.27.0

replace golang.org/x/sync => golang.org/x/sync v0.6.0

replace golang.org/x/xerrors => golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028

module example.com/log-nats-demo

go 1.26.6

toolchain go1.27.3

require github.com/taigrr/log-nats/v2 v2.1.1

require (
	github.com/klauspost/compress v1.20.3 // indirect
	github.com/nats-io/nats.go v1.53.1 // indirect
	github.com/nats-io/nkeys v0.4.16 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.58.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
)

replace github.com/taigrr/log-nats/v2 => ../

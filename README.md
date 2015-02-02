## Concurrent Downloader

Demonstrates concurrent downloads in [Go](https://golang.org/) via a worker pattern.  The number of
workers is parameterized.  Based on the [Go by Example](https://gobyexample.com) [Worker Pools example](https://gobyexample.com/worker-pools)

Generates a workload by generating random queries to google.  The number of parallel workers are spawned and start processing the jobs as they are sent in via [channels](https://gobyexample.com/channels).

[`main.go`](main.go) is the entrypoint

### Running

	go run main.go

> runs as script

### Build and Run

	go build
	./concurrent-downloader

### Install

	go install

> resulting `concurrent-downloader` binary in `$GOPATH/bin`
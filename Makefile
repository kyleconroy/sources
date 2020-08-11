example: $(wildcard *.go)
	go build -o example .

core/core.pb.go: core.proto
	protoc core.proto --go_out=./core

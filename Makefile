subspacePath=github.com/subspace-engine/subspace

# names and file locations are experimental
#executables will be at $GOPATH/bin

all:
	go install $(subspacePath)/engine/cmd/engine
	go install $(subspacePath)/engine/cmd/testCli


proto:
	protoc -I engine/ engine/message.proto --go_out=plugins=grpc:engine
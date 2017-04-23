all:
	go install ...

proto:
	protoc -I engine/ engine/message.proto --go_out=plugins=grpc:engine
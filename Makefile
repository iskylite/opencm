all: pb

.PHONY: pb

pb:
	protoc -I ./ --go_out=plugins=grpc:. pb/*.proto
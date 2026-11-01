PROTO_DIR   := proto
THIRD_PARTY := third_party
GEN_DIR     := gen/go
PROTO_FILES := $(shell find $(PROTO_DIR) -name '*.proto')

export PATH := $(PATH):$(shell go env GOPATH)/bin

.PHONY: proto-generate
proto-generate:
	@mkdir -p $(GEN_DIR)
	@protoc -I $(PROTO_DIR) -I $(THIRD_PARTY) -I /usr/include $(PROTO_FILES) \
		--go_out=$(GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative
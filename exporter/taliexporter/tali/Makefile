# Define directories
PROTO_DIR = proto
OUT_DIR = internal/gen_proto

# Define the .proto files and their generated targets
PROTO_FILES = $(wildcard $(PROTO_DIR)/*.proto)
GENERATED_FILES = $(patsubst $(PROTO_DIR)/%.proto, $(OUT_DIR)/%.pb.go, $(PROTO_FILES))

# Protoc command
PROTOC = protoc

# Default rule
all: $(GENERATED_FILES)

# Rule to generate Go files from .proto files
$(OUT_DIR)/%.pb.go: $(PROTO_DIR)/%.proto
	$(PROTOC) --proto_path=$(PROTO_DIR) \
	          --go_out=$(OUT_DIR) \
	          --go_opt=paths=source_relative $<

# Clean rule
clean:
	rm -rf $(OUT_DIR)/*.pb.go

.PHONY: all clean
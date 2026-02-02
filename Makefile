BUF_DIR := buf
THREAD_PROTO_INTERNAL := ../../protos/im/service/contact/v1
THREAD_PROTO_SHARED := ../../protos/im/domain/contact/v1

.PHONY: gen-contact

gen-contact:
	@echo "Generating thread protos"
	cd $(BUF_DIR)/ && go run github.com/bufbuild/buf/cmd/buf@latest generate \
		--template buf.gen.contact.yaml \
		--path $(THREAD_PROTO_INTERNAL) \
		--path $(THREAD_PROTO_SHARED)
	@echo "End of generating thread protos."
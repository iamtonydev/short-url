.PHONY: generate

generate:
	mkdir -p pkg/url_shortener_v1
	protoc --proto_path vendor.protogen --proto_path api/url_shortener_v1 --go_out=pkg/url_shortener_v1 --go_opt=paths=import \
          --go-grpc_out=pkg/url_shortener_v1 --go-grpc_opt=paths=import \
          --grpc-gateway_out=pkg/url_shortener_v1 \
          --grpc-gateway_opt=logtostderr=true \
          --grpc-gateway_opt=paths=import \
          --validate_out lang=go:pkg/url_shortener_v1 \
          api/url_shortener_v1/url_shortener.proto
	mv pkg/url_shortener_v1/github.com/iamtonydev/url-shortener/pkg/url_shortener_v1/* pkg/url_shortener_v1/
	rm -rf pkg/url_shortener_v1/github.com/

PHONY: vendor-proto
vendor-proto: .vendor-proto

PHONY: .vendor-proto
.vendor-proto:
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/github.com/envoyproxy ]; then \
			mkdir -p vendor.protogen/validate &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
			rm -rf vendor.protogen/protoc-gen-validate ;\
		fi

PHONY: run-migrations
run-migrations: .run-migrations

PHONY: .run-migrations
.run-migrations:
	sh ./migration.local.sh

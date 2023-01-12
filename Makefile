.PHONY: run-generators

run-generators: gen-api

gen-api:
	oapi-codegen --config codegen/configs/nuts-storage-api-server.yaml codegen/api-specs/nuts-storage-api-v1.yaml | gofmt > api/v1/generated.go

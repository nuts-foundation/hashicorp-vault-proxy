.PHONY: run-generators

run-generators: gen-api

gen-api:
	oapi-codegen --config codegen/configs/nuts-storage-api-server.yaml codegen/api-specs/nuts-storage-api-v1.yaml | gofmt > api/v1/generated.go

docker:
	docker build -t nutsfoundation/hashicorp-vault-proxy .

run:
	docker compose stop && docker compose rm -f
	docker compose up --wait
	docker compose exec -e VAULT_TOKEN=root vault vault secrets enable -version=1 -address=http://localhost:8200 kv

.PHONY: build start

run-generators: gen-api

gen-api:
	oapi-codegen --config codegen/configs/nuts-storage-api-server.yaml codegen/api-specs/nuts-storage-api-v1.yaml | gofmt > api/v1/generated.go

build:
	docker build -t nutsfoundation/hashicorp-vault-proxy .

start: stop start-vault
	docker compose up proxy --wait

stop:
	docker compose stop

start-vault:
	docker compose up vault --wait
	docker compose exec -e VAULT_TOKEN=root vault vault secrets enable -version=1 -address=http://localhost:8200 kv

reset:
	docker compose stop && docker compose rm -v -f

api-test: start
	docker compose up postman

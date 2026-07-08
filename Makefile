
.PHONY: test verify-dependency-security

test:
	go test ./...

verify-dependency-security:
	bash scripts/verify-dependency-security.sh

create_dummy_certs:
	./gen-certs.sh

start_server:
	go run main.go


run_client:
	go run main.go client

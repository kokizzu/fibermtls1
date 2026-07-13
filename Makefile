GO ?= go
GOVULNCHECK ?= govulncheck
CMD ?=

.PHONY: test verify-dependency-security create_dummy_certs create_dummy_mtls start_server run_client run vulncheck

test:
	$(GO) test ./...

verify-dependency-security:
	bash scripts/verify-dependency-security.sh

create_dummy_certs:
	./gen-certs.sh

create_dummy_mtls: create_dummy_certs

start_server:
	$(GO) run main.go


run_client:
	$(GO) run main.go client

vulncheck:
	$(GOVULNCHECK) ./...

run:
	@test -n "$(CMD)" || (echo "usage: make run CMD='go test ./...'" >&2; exit 2)
	$(CMD)

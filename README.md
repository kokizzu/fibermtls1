
# Simple Example how to create mTLS with Fiber

## Create mTLS certificate

generate both certificates

```bash
make create_dummy_mtls
```

## Start Server

start mTLS+HTTP server

```bash
make start_server
```

## Run Client

try to connect to server with mTLS

```bash
make run_client
```

## Maintenance checklist

- [x] Go runtime updated to 1.26.5.
- [x] Direct Go dependencies refreshed and module files tidied.
- [x] `make test` compiles and tests the mTLS example.
- [x] `make verify-dependency-security` and `make vulncheck` check dependency security.

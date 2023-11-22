
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

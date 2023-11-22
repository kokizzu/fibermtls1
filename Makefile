
create_dummy_certs:
	./gen-certs.sh

start_server:
	go run main.go


run_client:
	go run main.go client
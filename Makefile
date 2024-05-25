
run.debug:
	go run cmd/app/main.go -http-port=8080 -log-level=debug -db-url=postgres://user:pass@localhost:5432/postgres -pool=2


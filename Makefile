up:
	sudo docker compose -f docker-compose.yaml up --build -d

down:
	sudo docker compose -f docker-compose.yaml down

run.debug:
	go run cmd/app/main.go -http-port=8080 -log-level=debug -cache-url=localhost:6380

test:
	sudo docker compose -f docker-compose.test.yaml up --build --abort-on-container-exit
	sudo docker compose -f docker-compose.test.yaml down --volumes

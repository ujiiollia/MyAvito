up:
	docker compose up
start:
	go run cmd/app/main.go
down:
	docker rmi -f $(docker images -aq)

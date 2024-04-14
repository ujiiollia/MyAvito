up:
	docker compose up
run:
	go run cmd/app/main.go
down:
	docker rmi -f $(docker images -aq)

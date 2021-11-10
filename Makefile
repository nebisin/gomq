rabbit:
	docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.9-management

producer:
	go run cmd/producer/main.go

consumer:
	go run cmd/consumer/main.go

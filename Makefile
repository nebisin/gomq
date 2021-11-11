rabbit:
	docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.9-management

producer:
	go run cmd/producer/main.go $(TYPE) $(MESSAGE)

consumer:
	go run cmd/consumer/main.go info warning error

consumer-file:
	go run cmd/consumer/main.go warning error 2>> logs_from_rabbit.log
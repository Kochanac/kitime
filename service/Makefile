run:
	REDIS_HOST=localhost \
	KAFKA_HOST=kafka-cluster-kafka-bootstrap.infra:9092 \
	KAFKA_TOPIC=vobla-topic \
	CLICK_HOST=localhost \
	CLICK_USER=asd \
	CLICK_PASSWORD=asd \
	go run cmd/main.go

run local:
	REDIS_HOST=localhost \
	KAFKA_HOST=kafka:9092 \
	KAFKA_TOPIC=test-topic \
	CLICK_HOST=localhost:9000 \
	CLICK_USER=asd \
	CLICK_PASSWORD=asd \
	go run cmd/main.go

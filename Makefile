start:
	docker-compose up -d --build

stop:
	docker-compose down

stopv:
	docker-compose down --volumes

start-influx:
	docker-compose -f docker-compose.influxdb.yaml  up -d

stop-influx:
	docker-compose -f docker-compose.influxdb.yaml down

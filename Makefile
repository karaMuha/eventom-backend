setup:
	mkdir db-data; \
	mkdir db-data/postgres; \
	openssl genrsa -out private-key.pem 4096

start:
	docker-compose up -d
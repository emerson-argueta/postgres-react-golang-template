tests:
	go clean -testcache
	docker-compose down -v
	docker-compose up -d
	sleep 1
	godotenv -f ./test.env go test ./... -cover
	docker-compose down -v
test-database:
	docker-compose down -v
	docker-compose up
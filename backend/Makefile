tests:
	go clean -testcache
	docker-compose down -v
	docker-compose up -d
	sleep 1
	godotenv -f ./test.env go test ./... -cover
	docker-compose down -v
test-database:
	# build image only if it exists
	docker image inspect prgstack/community-goal-tracker-service:db-1.0 || (docker build . -t prgstack/community-goal-tracker-service:db-1.0)
	docker-compose down -v
	docker-compose up
	
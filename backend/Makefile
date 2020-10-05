tests:
	go clean -testcache
	cd ./.test && docker-compose down -v
	cd ./.test && docker-compose up -d
	sleep 1
	godotenv -f ./test.env go test ./... -cover
	cd ./.test && docker-compose down -v
test-database:
	cd ./.test && docker-compose down -v
	# build image only if it exists
	docker image inspect prgstack/community-goal-tracker-service:db-1.0 || (cd ./.test && docker build . -t prgstack/community-goal-tracker-service:db1.0)
	cd ./.test && docker-compose down -v
	cd ./.test && docker-compose up
	
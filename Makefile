
migrate-up:
	export $$(cat .env | xargs) && ./bin/migrate -database $${DB_CONNECTION_STR} -path migration up

migrate-down:
	export $$(cat .env | xargs) && ./bin/migrate -database $${DB_CONNECTION_STR} -path migration down

migrate-up-test:
	export $$(cat .env | xargs) && ./bin/migrate -database $${TEST_DB_CONNECTION_STR} -path migration up

migrate-down-test:
	export $$(cat .env | xargs) && ./bin/migrate -database $${TEST_DB_CONNECTION_STR} -path migration down
	
run:
	export $$(cat .env | xargs) && go run main.go

test: migrate-down-test migrate-up-test
	export $$(cat .env | xargs) && \
	go clean -testcache && go test -coverprofile cover.out github.com/syukur91/ischool-monitor/service -v

test-cover: migrate-down-test migrate-up-test
	export $$(cat .env | xargs) && \
	go clean -testcache && go test -coverprofile cover.out github.com/syukur91/ischool-monitor/service -v  && go tool cover -html=cover.out 

rebuild-db: migrate-down-test migrate-up-test

test-pkg: 
	go test github.com/syukur91/ischool-monitor/pkg/query -v

.PHONY: migrate-up migrate-down run test test-pkg
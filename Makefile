auth-docker:
	docker-compose build auth
	docker-compose up auth

dropdb:
	docker exec -it postgres12 dropdb table_mate

sqlc:
	sqlc generate --experimental

mock:
	mockgen -package mockrepo -destination db/mock/repo.go auth-service/db/repo Repo

gen:
	sqlc generate --experimental
	mockgen -package mockrepo -destination db/mock/repo.go auth-service/db/repo Repo

test:
	go test -v -cover ./...

run:
	go run .

test-repo:
	go test -v ./db/repo -run Test

.PHONY: auth-docker sqlc mock test run test-repo gen

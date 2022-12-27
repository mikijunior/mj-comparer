run:
	go run ./cmd/apiserver/

build:
	go build ./cmd/apiserver/

migrate-up:
	migrate -path ./migrations/ -database "mysql://$(u):$p@tcp($(host):3306)/$(name)" -verbose up

migrate-down:
	migrate -path ./migrations/ -database "mysql://$(u):$p@tcp($(host):3306)/$(name)" -verbose down

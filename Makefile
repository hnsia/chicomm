include .env

mysql-up:
	docker run --name chicomm-mysql -p 3305:3306 -e MYSQL_ROOT_PASSWORD=password -d mysql:8.4

mysql-create-db:
	docker exec -i chicomm-mysql mysql -uroot -ppassword <<< "CREATE DATABASE chicomm;"

mysql-start:
	docker start chicomm-mysql

mysql-stop:
	docker stop chicomm-mysql

mysql-remove:
	docker stop chicomm-mysql
	docker rm chicomm-mysql

migrate-create:
	docker run -it --rm --network host --volume "${PWD}/db:/db" migrate/migrate:v4.17.0 create -ext sql -dir /db/migrations $(filename)

n?=1
migrate-up:
	docker run -it --rm --network host --volume "$(PWD)/db:/db" migrate/migrate:v4.17.0 -path=/db/migrations -database "mysql://root:password@tcp(localhost:3305)/chicomm" up $(n)

migrate-down:
	docker run -it --rm --network host --volume "$(PWD)/db:/db" migrate/migrate:v4.17.0 -path=/db/migrations -database "mysql://root:password@tcp(localhost:3305)/chicomm" down $(n)

grpc-codegen:
	protoc --proto_path=chicomm-grpc/pb --go_out=chicomm-grpc/pb --go_opt=paths=source_relative \
    --go-grpc_out=chicomm-grpc/pb --go-grpc_opt=paths=source_relative \
    chicomm-grpc/pb/api.proto

run-grpc:
	go run cmd/chicomm-grpc/main.go

run-api:
	go run cmd/chicomm-api/main.go

run-notification:
	ADMIN_EMAIL=$(ADMIN_EMAIL) ADMIN_PASS=$(ADMIN_PASS) go run cmd/chicomm-notification/main.go

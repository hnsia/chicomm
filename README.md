# Chicomm
This is a e-commerce dummy Golang project to explore chi-router, grpc, and sqlx.

# Setup mysql

1. Run `docker pull mysql:8.4`
2. Run `docker run --name chicomm-mysql -p 3305:3306 -e MYSQL_ROOT_PASSWORD=password -d mysql:8.4` or `make mysql-up`
3. Run `docker exec -i chicomm-mysql mysql -uroot -ppassword <<< "CREATE DATABASE chicomm;"` or `make mysql-create-db`

# Run mysql container (if already setup)
1. Run `make mysql-start`

# Run migrate

1. Run `docker run -it --rm --network host --volume "$(pwd)/db:/db" migrate/migrate:v4.17.0 create -ext sql -dir /db/migrations init_schema` or `make migrate-create filename="init_schema"` to create migration scripts
2. Run `docker run -it --rm --network host --volume "$(pwd)/db:/db" migrate/migrate:v4.17.0 -path=/db/migrations -database "mysql://root:password@tcp(localhost:3305)/chicomm" up` or `make migrate-up n=1` to migrate up
3. Run `docker run -it --rm --network host --volume "$(pwd)/db:/db" migrate/migrate:v4.17.0 -path=/db/migrations -database "mysql://root:password@tcp(localhost:3305)/chicomm" down` or `make migrate-down n=1` to migrate down

# To generate grpc code
1. Run 
```
protoc --proto_path=chicomm-grpc/pb --go_out=chicomm-grpc/pb --go_opt=paths=source_relative \
    --go-grpc_out=chicomm-grpc/pb --go-grpc_opt=paths=source_relative \
    chicomm-grpc/pb/api.proto
```
or `make grpc-codegen`

# Environment variables
1. make a copy of .env.example and save it as .env
2. fill in the required environment variables

# Run services
1. Run `make run-grpc` to start grpc server
2. Run `make run-api` to start api server
3. Run `make run-notification` to start notification server

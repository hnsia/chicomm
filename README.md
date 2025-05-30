# Chicomm

# Setup mysql

1. `docker pull mysql:8.4`
2. `docker run --name chicomm-mysql -p 3305:3306 -e MYSQL_ROOT_PASSWORD=password -d mysql:8.4`
3. `docker exec -i chicomm-mysql mysql -uroot -ppassword <<< "CREATE DATABASE chicomm;"`

# Run migrate

1. Run `docker run -it --rm --network host --volume "$(pwd)/db:/db" migrate/migrate:v4.17.0 create -ext sql -dir /db/migrations init_schema` to create migration scripts
2. Run `docker run -it --rm --network host --volume "$(pwd)/db:/db" migrate/migrate:v4.17.0 -path=/db/migrations -database "mysql://root:password@tcp(localhost:3305)/chicomm" up` to migrate up
3. Run `docker run -it --rm --network host --volume "$(pwd)/db:/db" migrate/migrate:v4.17.0 -path=/db/migrations -database "mysql://root:password@tcp(localhost:3305)/chicomm" down` to migrate down

# To generate grpc code
1. Run 
```
protoc --proto_path=chicomm-grpc/pb --go_out=chicomm-grpc/pb --go_opt=paths=source_relative \
    --go-grpc_out=chicomm-grpc/pb --go-grpc_opt=paths=source_relative \
    chicomm-grpc/pb/api.proto
```
mysql-up:
	docker run --name chicomm-mysql -p 3305:3306 -e MYSQL_ROOT_PASSWORD=password -d mysql:8.4

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
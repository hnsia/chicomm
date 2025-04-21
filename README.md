# Chicomm

# Setup mysql

1. `docker pull mysql:8.4`
2. `docker run --name chicomm-mysql -p 3305:3306 -e MYSQL_ROOT_PASSWORD=password -d mysql:8.4`

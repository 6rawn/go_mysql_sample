# go_mysql_sample

## Prepare a sample database

```bash
$ git clone https://github.com/6rawn/mysql-replication-sample
$ cd mysql-replication-sample
$ docker compose up -d
$ ./run.sh

$ curl -L -O https://github.com/catatsuy/private-isu/releases/download/img/dump.sql.bz2
$ bunzip2 dump.sql.bz2
$ mv dump.sql scripts
$ docker compose exec -it node-1 bash
$ mysql -u root -pmysql < scripts/dump.sql
$ exit
```

## Usage

```
$ DB_USER=root DB_PASSWORD=mysql go run main.go
```

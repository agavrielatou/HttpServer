docker run --net="host" --rm -P --name pg_test eg_postgresql
#docker run --rm -t -i --link pg_test:pg eg_postgresql bash
postgres@7ef98b1b7243:/$ psql -h $PG_PORT_5432_TCP_ADDR -p $PG_PORT_5432_TCP_PORT -d docker -U docker --password

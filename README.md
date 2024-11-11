docker run -e POSTGRES_PASSWORD=password -e POSTGRES_USER=dbuser -e POSTGRES_DB=test -p 5432:5432 postgres

migrate --path database/migration/ -database "postgresql://dbuser:password@localhost:5432/test?sslmode=disable" -verbose up

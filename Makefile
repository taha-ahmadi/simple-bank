postgres:
	docker run --name simple-bank-postgres -p 5432:5432 -e POSTGRES_USER=taha -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it simple-bank-postgres createdb --username=taha --owner=taha simple_bank

dropdb:
	docker exec -it simple-bank-postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://taha:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://taha:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc
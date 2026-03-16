migrateup:
	migrate -path migrations -database "postgresql://myuser:mypassword@localhost:5432/mydb?sslmode=disable" -verbose up

migratedown:
	migrate -path migrations -database "postgresql://myuser:mypassword@localhost:5432/mydb?sslmode=disable" -verbose down

migrateup1:
	migrate -path migrations -database "postgresql://myuser:mypassword@localhost:5432/mydb?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path migrations -database "postgresql://myuser:mypassword@localhost:5432/mydb?sslmode=disable" -verbose down 1

new_migration:
	migrate create -ext sql -dir migrations -seq mydb

createdb:
	docker exec postgres-db createdb --username=myuser --owner=myuser mydb

dropdb:
	docker exec postgres-db dropdb --username=myuser mydb

server:
	go run main.go

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: migrateup migratedown migrateup1 migratedown1 createdb dropdb server test new_migration
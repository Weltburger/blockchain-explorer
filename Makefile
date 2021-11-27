PWD = $(shell pwd)
MPATH = $(PWD)/internal/auth/migrations

migrate-create:
	@echo "---Creating migration files---"
	migrate create -ext sql -dir $(MPATH) -digits 3 -seq $(NAME)

postgres:
	docker run --rm --name postgres_block-explorer --network host -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres

createdb:
	docker exec -it postgres_block-explorer createdb --username=root --owner=root blockchain_explorer

dropdb:
	docker exec -it postgres_block-explorer dropdb blockchain_explorer

migrate-up:
	migrate -path $(MPATH) -database "postgresql://root:password@localhost:5432/blockchain_explorer?sslmode=disable" -verbose up

migrate-down:
	migrate -path $(MPATH) -database "postgresql://root:password@localhost:5432/blockchain_explorer?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb migrate-up migrate-down

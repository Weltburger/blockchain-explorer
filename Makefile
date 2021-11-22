.PHONY: migrate-create migrate-up migrate-down migrate-force dock-postgres-start

PWD = $(shell pwd)
ACCTPATH = $(PWD)
MPATH = $(ACCTPATH)/cmd/migration
PORT = 5432

# Default number of migrations to execute up or down
N = 1

# Create keypair should be in your file below

migrate-create:
	@echo "---Creating migration files---"
	migrate create -ext sql -dir $(MPATH) -digits 5 -seq $(NAME)

dock-postgres-start:
	docker run -it --rm -d --network host --name postgres_block-explorer -e POSTGRES_PASSWORD=password -e POSTGRES_USER=postgres postgres

migrate-up:
	migrate -source file://$(MPATH) -database postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable up $(N)

migrate-down:
	migrate -source file://$(MPATH) -database postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable down $(N)

migrate-force:
	migrate -source file://$(MPATH) -database postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable force $(VERSION)
PWD = $(shell pwd)
MPATH = $(PWD)/internal/auth/migrations

create-keypair:
	@echo "Creating an rsa 256 key pair"
	openssl genpkey -algorithm RSA -out $(PWD)/rsa_private_$(ENV).pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in $(PWD)/rsa_private_$(ENV).pem -pubout -out $(PWD)/rsa_public_$(ENV).pem

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

.PHONY: postgres createdb dropdb migrate-up migrate-down migrate-create create-keypair


#docker run -d --name pg_db --hostname pgdb --network blockchainnet -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -e PGDATA=/var/lib/postgresql/data/pgdata -e POSTGRES_DB=blockchain_explorer -v pgvolume:/var/lib/postgresql/data postgres

#docker run -it --rm --network blockchainnet --link ch_db:clickhouse-server yandex/clickhouse-client --host clickhouse-server

#docker run -d --name ch_db --hostname chdb --network blockchainnet -p 8123:8123 -p 9100:9100 --ulimit nofile=262144:262144 -v chvolume:/var/lib/clickhouse yandex/clickhouse-server

#docker run -it --rm --network blockchain-explorer_blocknet --name explorer1 -p 1323:1323 -e HTTP_PORT=:1323 -e SIGNING_KEY=strongPassword -e TOKEN_TTL=86400 -e STEP=1000 -e NODE=https://testnet-tezos.giganode.io/chains/main/blocks/ -e POSTGRES_HOST=pgdb -e POSTGRES_PORT=5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=mysecretpassword -e POSTGRES_DB=blockchain_explorer -e PGDATA=/var/lib/postgresql/data/ -e CLICKHOUSE_HOST=chdb -e CLICKHOUSE_PORT=8123 -e CLICKHOUSE_DB=blocks -e CLICKHOUSE_DEBUG=true -e TOTAL_WORKER=10 -e CRAWLER_START_POS=678500 -v /home/phandorin/dev/blockchain-explorer/internal/storage/migration:/migrationsClickhouse -v /home/phandorin/dev/blockchain-explorer/internal/auth/migrations:/migrationsPostgres explorer:v0.9

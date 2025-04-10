VERSION=0.0.1
BUILD=`date +%FT%T%z`
GIT_HASH=`git rev-parse HEAD`
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD) -X main.GitHash=$(GIT_HASH)"
# sql-migrate
CMD_RM=rm -rf
APP_MIGRATE_CFG=resources/migrations/$*/app/dbconfig.yml
SCHEMA_MIGRATE_CFG=resources/migrations/$*/schema/dbconfig.yml
MODEL_DIR=pkg/$*/models
CMD_BOILER=sqlboiler -c resources/migrations/$*/sqlboiler.toml -o $(MODEL_DIR)
CMD_MIGRATE=sql-migrate

.PHONY: test
test:
	go test -v ./...

LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD) -X main.GitHash=$(GIT_HASH)"

.PHONY: run-%
run-%:
	go run -v $(LDFLAGS) cmd/$*/main.go

.PHONY: docker-up
docker-up:
	@docker network ls | awk '{print $2}' | grep springboard > /dev/null 2>&1 || docker network create tr
	docker-compose -f docker-compose.yaml up --build -d --remove-orphans

.PHONY: docker-down
docker-down:
	docker-compose -f docker-compose.yaml down

.PHONY: logs-%
logs-%:
	@docker ps --format '{{.Names}}' | grep $* | xargs docker logs -f

.PHONY: new-%
new-%:
	mkdir cmd/$* && cd cmd/$* && touch main.go
	mkdir internal/$* && cd internal/$* && mkdir api config app repository
	cd internal/$*/api && touch api.go routes.go
	cd internal/$*/config && touch config.go config_toml.go

##
## database migration
##
.PHONY: db-migrate-status-%
db-migrate-status-%:
	$(CMD_MIGRATE) status -config=$(SCHEMA_MIGRATE_CFG)
	$(CMD_MIGRATE) status -config=$(APP_MIGRATE_CFG)

.PHONY: db-migrate-up-%
db-migrate-up-%:
	$(CMD_MIGRATE) up -config=$(APP_MIGRATE_CFG)

.PHONY: db-migrate-down-%
db-migrate-down-%:
	$(CMD_MIGRATE) down -config=$(APP_MIGRATE_CFG)

.PHONY: db-clean-models-%
db-clean-models-%:
	$(CMD_BOILER) mysql --wipe
	$(CMD_RM) $(MODEL_DIR)

.PHONY: db-generate-models-%
db-generate-models-%:
	$(CMD_BOILER) mysql --wipe
	$(CMD_RM) $(MODEL_DIR)
	$(CMD_BOILER) mysql --no-hooks --no-tests

##
## database schema migrations
##
.PHONY: db-reset-%
db-reset-%:
	$(CMD_MIGRATE) down -config=$(SCHEMA_MIGRATE_CFG) -limit=0
	make db-init

.PHONY: db-init-%
db-init-%:
	$(CMD_MIGRATE) up -config=$(SCHEMA_MIGRATE_CFG)

include .env

GOBASE=$(shell pwd)

MIGRATOR_DIR=${GOBASE}/migrations/

refers-run:
	@echo "  >  Building search service ...."
	go build -o ./cmd/refers/service ./cmd/refers
	@echo "  >  Search run ...."
	go run ./cmd/refers/main.go serve


# миграция баазы
mig-up:
	go run main.go up -d ${MIGRATOR_DIR} -s $(schema)
mig-down:
	go run main.go down	-d ${MIGRATOR_DIR} -s $(schema)
mig-c:
	go run main.go create -n $(name) -d ${MIGRATOR_DIR} -s $(schema)

mig-r:
	go run main.go reset  -d ${MIGRATOR_DIR} -s $(schema)

mig-sc:
	go run main.go schema -d ${MIGRATOR_DIR} -s $(schema)

clear: db-down db-rm db-up


#  make mig-c schema=public name=users  создание миграции
#  накат миграции
#  make mig-up schema=public name=users  создание миграции
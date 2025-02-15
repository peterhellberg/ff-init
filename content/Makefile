APP_ID=ff-app-id
AUTHOR_ID=ff-author-id
ARCHIVE=${AUTHOR_ID}.${APP_ID}.zip
HOSTNAME=localhost
APP_PATH=ff/${APP_ID}
APP_URL=https://${HOSTNAME}/${APP_PATH}
SERVER_ROOT=/tmp
SERVER_PATH=${SERVER_ROOT}/${APP_PATH}

.PHONY: all
all: spy

.PHONY:
spy:
	zig build spy

.PHONY:
run:
	zig build run

.PHONY: build
build:
	zig build

.PHONY: export
export: build
	@firefly_cli export

.PHONY: deploy
deploy: export
	@ssh ${HOSTNAME} 'mkdir -p ${SERVER_PATH}/src ${SERVER_PATH}/rom'
	@scp -q *.zig ${HOSTNAME}:${SERVER_PATH}/
	@scp -q *.zon ${HOSTNAME}:${SERVER_PATH}/
	@scp -q *.toml ${HOSTNAME}:${SERVER_PATH}/
	@scp -q *.fff ${HOSTNAME}:${SERVER_PATH}/
	@scp -q *.md ${HOSTNAME}:${SERVER_PATH}/
	@scp -q src/*.zig ${HOSTNAME}:${SERVER_PATH}/src/
	@echo "✔ Updated ${APP_ID} on ${APP_URL}"
	@scp -q ${ARCHIVE} ${HOSTNAME}:${SERVER_PATH}/rom/${ARCHIVE}
	@echo "✔ Archive ${APP_URL}/rom/${ARCHIVE}"

all: build

update:
	git submodule update --init --remote --recursive
	git submodule sync --recursive

build: update
	mkdir -p build/bin
	env go build -o build/bin/app main.go

release: build
	rm bin/start
	cp build/bin/app bin/start
	chmod a+x bin/start

deploy: release
	bash scripts/Deploy.sh

dump-schema:
	mkdir -p cloud/init
	sudo -u postgres pg_dump --dbname academy --schema-only >cloud/init/academy.psql

docker-compose-build:
	sudo env DOCKER_BUILDKIT=1 docker-compose -f compose.yml -p academy build

docker-compose-up: docker-compose-build
	sudo docker-compose up

docker-build:
	sudo env DOCKER_BUILDKIT=1 docker build -f Dockerfile -t academy .

docker-down:
	docker-compose down

docker-up: docker-compose-up

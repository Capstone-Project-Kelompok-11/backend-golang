all: build

update:
	git submodule update --init --recursive

build: update
	mkdir -p bin
	env go build -o bin/app main.go

release: build
	chmod a+x bin/start

deploy: release
	bash scripts/Deploy.sh

dump-schema:
	mkdir -p cloud/init
	sudo -u postgres pg_dump --dbname academy --schema-only >cloud/init/academy.psql

docker-compose-build:
	sudo env DOCKER_BUILDKIT=1 docker compose -f docker-compose.yaml -p academy build

docker-compose-up: docker-compose-build
	sudo docker compose up

docker-build:
	sudo env DOCKER_BUILDKIT=1 docker build -f Dockerfile -t academy .

docker-down:
	docker-compose down

docker-up: docker-compose-up

remove-assets-caches:
	rm -rvf assets/public/caches/*

test: remove-assets-caches
	exec bin/start . &>/dev/null
	go run test/main/wait_tcp_open.go
	go test -v test/unit_test.go
	bash scripts/ForceStop.sh

prune:
	docker buildx prune -f
	docker container prune -f
	docker image prune -f
	docker network prune -f
	docker system prune -f

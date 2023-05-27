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

dump-schema:
	mkdir -p cloud/init
	sudo -u postgres pg_dump --dbname academy --schema-only >cloud/init/academy.psql

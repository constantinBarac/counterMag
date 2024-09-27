# Change these variables as necessary.
main_package_path = ./cmd/countermag
binary_name = countermag

.PHONY: test/unit
test/unit:
	go test -v -race ./...

.PHONY: test/load
test/load: build
	scripts/load.sh

.PHONY: test/replication
test/replication: build
	python3 tests/replication/replication.py

.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

## build: build the application
.PHONY: build
build:
	go build -o=./bin/${binary_name} ${main_package_path}

## run: run the  application
.PHONY: run
run: build
	./bin/${binary_name}

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "./bin/${binary_name}" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"

.PHONY: run/replicated
run/replicated: build 
	./scripts/spinup_fg.sh
	./scripts/teardown.sh

.PHONY: clean
clean:
	./scripts/teardown.sh
	rm -r out/
	rm -r bin/


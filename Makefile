include .envrc

.PHONY: help
## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: dev
## dev: start go-air, enable hot reload
dev:
	watchexec -r -e go --clear --wrap-process session -- "go run main.go"

.PHONY: dev/init
dev/init:
	watchexec -r -e go --clear --wrap-process session -- "go run main.go init"


## tidy: format all .go files and tidy module dependencies
.PHONY: tidy
tidy:
	@echo 'Formatting .go files...'
	go fmt ./...
	@echo 'Tidying module dependencies...'
	go mod tidy

## audit: run quality control checks
.PHONY: audit
audit:
	@echo 'Checking module dependencies'
	go mod tidy -diff
	go mod verify
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...


.PHONY: run
## run: run the app
run:
	go run main.go $(filter-out $@,$(MAKECMDGOALS))

%:
	@:


build:
	@go build -o bin/go-migrations

run: build
	@./bin/go-migrations

renamefs:
	@./sh-testing-cmds/.bs-renamefs.sh
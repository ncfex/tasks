build:
	@go build -o tasks.out cmd/tasks/main.go

run: build
	./tasks.out $(ARGS)
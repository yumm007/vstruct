all:
	go build -o bin/vstruct main.go
	go generate ./...

all:
	#find . -name "*_vstruct.go" | xargs -I{} rm {}
	go build -o bin/vstruct main.go
	go generate ./...

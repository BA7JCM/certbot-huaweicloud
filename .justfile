set shell := ["powershell.exe", "-c"]

build:
	go build -o ./bin/ ./...

auth:
	go run main.go -use auth

alias cp := cleanup
cleanup:
	go run main.go -use cleanup


set shell := ["powershell.exe", "-c"]

build:
	go build -o ./bin/ ./...

run:
	go run main.go
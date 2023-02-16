BIN_OUTPUT = bin/

all:
	#bash GOBIN=${PWD}/$(BIN_OUTPUT) && go install github.com/cosmtrek/air@latest
	#go build -o bin/web_backend ./cmd/web
	./bin/air --build.cmd "go build -o bin/web_backend cmd/web/main.go" --build.bin "./bin/web_backend"
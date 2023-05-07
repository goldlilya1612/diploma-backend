BIN_OUTPUT = bin/

run:
	go run ./cmd/web/main.go

build:
	go build -o bin/web_backend ./cmd/web

air:
	#bash GOBIN=${PWD}/$(BIN_OUTPUT) && go install github.com/cosmtrek/air@latest
	./bin/air --build.cmd "go build -o bin/web_backend cmd/web/main.go" --build.bin "./bin/web_backend"

docker-build:
	sudo docker build -t web_backend -f ./deployments/docker/Dockerfile .

docker-push:
	sudo docker tag web_backend sergencio/web_backend
	sudo docker push sergencio/web_backend

docker-compose-up:
	sudo docker-compose -f ./deployments/docker/docker-compose.yaml pull && \
	sudo docker-compose -f ./deployments/docker/docker-compose.yaml up

docker-compose-down-all:
	sudo docker-compose -f ./deployments/docker/docker-compose.yaml down -v --rmi all

docker-compose-start:
	sudo docker-compose -f ./deployments/docker/docker-compose.yaml start

docker-compose-stop:
	sudo docker-compose -f ./deployments/docker/docker-compose.yaml stop

clean:
	rm -rf ./bin ./tmp
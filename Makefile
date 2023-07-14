tidy: # @HELP go mod tidy and update vendor
tidy: 
	go mod tidy
	go mod vendor

test: # @HELP go run application test
test:
	go test ./... -count=1 -v

build: # @HELP build application binary and place in bin directory
build:
	go build \
		-o bin/kafka \
		./cmd/kafka

docker: # @HELP build application and run in docker
docker:
	docker build --rm -t kafka-server .
	docker-compose up --build

kafka-cluster: # @HELP build kafka cluster
kafka-cluster:
	docker-compose -f docker-compose-kafka-cluster.yml up --build

app:
app:
	docker-compose -f docker-compose.yml --build -f docker-compose-kafka-cluster.yml up --build

dev:
	docker-compose -f docker-compose-dev.yml up --build
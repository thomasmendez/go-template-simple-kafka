FROM golang:1.19.2-buster as build

WORKDIR /app/

COPY go.* /app/

COPY . .

RUN go mod download
RUN go mod verify

RUN make build

RUN go install github.com/githubnemo/CompileDaemon@latest
ENTRYPOINT CompileDaemon -log-prefix=false -exclude '*_test.go' -exclude-dir=.git -verbose -polling -build="go build -o /go/bin/ ./..." -command="kafka"
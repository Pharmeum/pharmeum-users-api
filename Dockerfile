FROM golang:1.12.1-stretch

WORKDIR $GOPATH/src/github.com/Pharmeum/pharmeum-users-api/

COPY . .

RUN go build -o users-api -v ./cmd/main.go
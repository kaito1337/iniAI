FROM golang:1.21

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go install github.com/cosmtrek/air@v1.49.0

RUN go mod download

COPY . .

RUN go build -o /server ./cmd/main/main.go

EXPOSE 7070

CMD ["/server"]
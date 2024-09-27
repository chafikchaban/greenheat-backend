FROM golang:1.23.1

WORKDIR /app

COPY go.mod .

COPY main.go .

RUN go build -o bin .

ENTRYPOINT ["/app/bin"]

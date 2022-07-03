FROM golang:1.18.3-buster

WORKDIR /app

COPY . .
RUN go mod download
WORKDIR /app/cmd
RUN go build -o main .

CMD ["./main"]
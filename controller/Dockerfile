FROM golang:1.18.3-buster

WORKDIR /go/src/github.com/dimaglushkov/dkvs/controller

COPY ../cmd/controller .

# Copy the source from the current directory to the Working Directory inside the container
COPY ../cmd/controller .

RUN go build -o main .

CMD ["./main"]
# Build stage
FROM golang AS builder

WORKDIR /todo-app

COPY . .

RUN go mod download
RUN go build -o todo todo.go


# Run stage
FROM ubuntu

WORKDIR /todo

COPY --from=builder /todo-app/todo .
RUN mkdir database

VOLUME ["/var/db"]
EXPOSE 8080

CMD ["./todo"]
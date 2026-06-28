FROM golang:1.25 AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler .

FROM ubuntu:latest

WORKDIR /app
COPY --from=builder /app/scheduler .
COPY web ./web

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/data/scheduler.db
# TODO_PASSWORD специально не реализую здесь, в виду того что это секрет

CMD ["./scheduler"]
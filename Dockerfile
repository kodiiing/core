FROM golang:1.18.2-bullseye AS builder
WORKDIR /app
COPY . .
RUN go build .

FROM debian:bullseye
WORKDIR /app
COPY --from=builder /app/kodiiing .
EXPOSE ${PORT}

CMD ["./kodiiing"]

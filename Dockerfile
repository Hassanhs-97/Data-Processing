FROM docker.iranrepo.ir/golang:1.21-alpine as builder

WORKDIR /app
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -v -o main .



FROM docker.iranrepo.ir/golang:1.20-alpine
WORKDIR /app
COPY --from=builder /app/main ./main

CMD ["./main"]
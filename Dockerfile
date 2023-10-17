FROM golang:1.20-alpine as builder

WORKDIR /app
COPY . /app

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -v -o app .

FROM golang:1.20-alpine

COPY --from=builder /app/app /app

ENTRYPOINT ["/app"]
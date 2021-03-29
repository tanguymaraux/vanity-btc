ARG VERSION=latest
FROM golang:${VERSION} as builder

WORKDIR /app
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . /app
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s"
FROM alpine
COPY --from=builder /app/vanity-btc /app/vanity-btc
ENTRYPOINT ["/app/vanity-btc"]
FROM golang:1.22-alpine as go-builder

WORKDIR /src/api

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

COPY api api/
COPY go.mod .
COPY go.sum .

RUN go build -o bin/mlp-api ./api/main.go

# Clean image with mlp-api binary
FROM alpine:3.16

COPY --from=go-builder /src/api/bin/mlp-api /usr/bin/mlp
COPY db-migrations ./db-migrations

ENTRYPOINT ["sh", "-c", "mlp \"$@\"", "--"]
CMD ["serve"]

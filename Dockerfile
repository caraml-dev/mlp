# ============================================================
# Build stage 1: Build UI
# ============================================================
FROM node:16-alpine as node-builder
WORKDIR /src/ui
COPY ui .
RUN yarn
RUN yarn lib build
RUN yarn app build

# ============================================================
# Build stage 2: Build API
# ============================================================
FROM golang:1.18-alpine as go-builder
WORKDIR /src/api
COPY api api/
COPY go.mod .
COPY go.sum .
COPY db-migrations ./db-migrations
RUN go build -o bin/mlp-api ./api/cmd/main.go

# ============================================================
# Build stage 3: Run the app
# ============================================================
FROM alpine:3.12
COPY --from=node-builder /src/ui/build ./ui/build
COPY --from=go-builder /src/api/bin/mlp-api /usr/bin/mlp
COPY --from=go-builder /src/api/db-migrations ./db-migrations

CMD ["mlp"]

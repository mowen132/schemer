# Copyright (c) 2025 Mark Owen
# Licensed under the MIT License. See LICENSE file in the project root for details.

FROM golang:1.24 AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o schemer ./cmd/schemer

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/schemer /usr/bin/schemer
RUN chmod +x /usr/bin/schemer

ENTRYPOINT ["schemer"]

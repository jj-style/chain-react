FROM golang:1.25-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN apk add --no-cache just
RUN just build


FROM scratch
WORKDIR /app

COPY --from=builder /src/build .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENV GIN_MODE=release
EXPOSE 8000
USER 1000:1000
ENTRYPOINT [ "./main", "server" ]

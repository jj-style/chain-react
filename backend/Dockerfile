FROM golang:1.21-alpine as builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN apk add --no-cache just
RUN just build


FROM scratch
WORKDIR /app

COPY --from=builder /src/build .
EXPOSE 8000
USER 1000:1000
ENTRYPOINT [ "./main", "server" ]

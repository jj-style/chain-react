FROM golang:1.20 as builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build


FROM bitnami/minideb
WORKDIR /app

COPY --from=builder /src/build .
EXPOSE 8000
USER 1000:1000
ENTRYPOINT [ "./main", "server" ]
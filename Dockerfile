FROM golang:1.21 AS builder
WORKDIR /app/

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o microservice /app/main.go

FROM alpine:latest
COPY --from=builder /app/microservice /microservice
COPY --from=builder /app/db /db

ENV GEOIP_DB_PATH=/db/GeoLite2-Country.mmdb
ENV CUSTOMERS_DB_PATH=/db/customers.db

ENV PORT=8080
EXPOSE 8080

CMD ["/microservice"]

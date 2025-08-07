FROM golang:1.20 AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o venue-searchpos main.go

FROM gcr.io/distroless/base

COPY --from=builder /app/venue-searchpos /usr/local/bin/venue-searchpos-service

CMD ["venue-searchpos-service"] 

EXPOSE 8080
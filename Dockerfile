FROM golang:alpine AS builder
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -o authentication ./...

FROM alpine:latest
RUN mkdir app
WORKDIR /app
COPY --from=builder /app/authentication .
CMD ["./authentication"]

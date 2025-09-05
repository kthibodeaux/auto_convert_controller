FROM golang:latest AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/server .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/server /app/server
CMD [ "./server" ]

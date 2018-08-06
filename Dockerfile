FROM golang:1.10 as builder
WORKDIR /go/src/github.com/matthewdale/matthewrdale.com/
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /go/src/github.com/matthewdale/matthewrdale.com/app .
COPY --from=builder /go/src/github.com/matthewdale/matthewrdale.com/public/ ./public
CMD ["./app"]

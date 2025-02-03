FROM golang:1.20-alpine AS builder
WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY cmd ./cmd
COPY templates ./templates
COPY static ./static


WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/app .


FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/app .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./app"]

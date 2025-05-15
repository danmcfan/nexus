FROM golang:1.24.3-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags=production -o nexus /app/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/nexus .
EXPOSE 80
ENV GIN_MODE=release
CMD ["./nexus"]
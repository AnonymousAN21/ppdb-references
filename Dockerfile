# ---------- Build Stage ----------
    FROM golang:1.24-alpine AS builder

    WORKDIR /app
    
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy everything (api/, config/, etc.)
    COPY . .
    
    # 👇 Build the Go binary using main.go inside /api
    RUN CGO_ENABLED=0 GOOS=linux go build -o main ./api
    
    # ---------- Run Stage ----------
    FROM alpine:3.21
    
    RUN apk --no-cache add ca-certificates
    
    WORKDIR /root/
    COPY --from=builder /app/main .
    
    EXPOSE 8080
    CMD ["./main"]
    
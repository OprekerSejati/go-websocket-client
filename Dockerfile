# Menggunakan image Go resmi sebagai base image  
FROM golang:1.20 AS builder  
  
# Set working directory  
WORKDIR /app  
  
# Menyalin go.mod dan go.sum  
COPY go.mod go.sum ./  
  
# Mengunduh dependensi  
RUN go mod download  
  
# Menyalin kode sumber  
COPY . .  
  
# Membangun aplikasi untuk Linux amd64  
RUN GOOS=linux GOARCH=amd64 go build -o client .  
  
# Menggunakan image minimal untuk menjalankan aplikasi  
FROM alpine:latest  
  
# Menyalin binary dari builder  
COPY --from=builder /app/client .  
  
# Menjalankan aplikasi  
CMD ["./client"]
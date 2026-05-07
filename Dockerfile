# ==========================================
# Tahap 1: Builder (Pabrik Perakitan)
# ==========================================
FROM golang:1.26.1-alpine AS builder

# Set folder kerja di dalam kontainer
WORKDIR /app

# Salin file konfigurasi modul terlebih dahulu (untuk caching)
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh kode aplikasi (termasuk folder templates)
COPY . .

# Melakukan kompilasi program Golang menjadi file tunggal bernama 'main'
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ==========================================
# Tahap 2: Production (Hasil Akhir)
# ==========================================
FROM alpine:latest

WORKDIR /app

# Hanya salin hasil kompilasi (mesin API) dari Tahap 1
COPY --from=builder /app/main .

# Salin folder templates agar HTML bisa ditampilkan di browser
COPY --from=builder /app/templates ./templates

# Buka port 8080
EXPOSE 8080

# Perintah wajib saat kontainer dijalankan
CMD ["./main"]
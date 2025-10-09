# ===== STAGE 1: build (Go 1.23 + bcg729) =====
FROM golang:1.23-alpine AS build

LABEL maintainer="Federico Pereira <fpereira@cnsoluciones.com>"

# Herramientas nativas para compilar bcg729 en Alpine
RUN apk add --no-cache \
    build-base \
    git \
    cmake \
    pkgconfig \
    ca-certificates

# Compilar e instalar bcg729 en /usr/local
RUN git clone https://github.com/BelledonneCommunications/bcg729 /tmp/bcg729 \
 && cd /tmp/bcg729 \
 && cmake -S . -B build -DBUILD_SHARED_LIBS=ON \
 && cmake --build build --target install

# App Go
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY transcoding.go ./

# CGO habilitado para enlazar con libbcg729.so (rpath ya en #cgo LDFLAGS)
ENV CGO_ENABLED=1
RUN go build -o /app/transcoding

# ===== STAGE 2: runtime =====
FROM alpine:latest AS runtime

LABEL maintainer="Federico Pereira <fpereira@cnsoluciones.com>"

# Instalar dependencias mínimas para runtime
RUN apk add --no-cache \
    ca-certificates \
    libc6-compat

# Copiamos la .so y el binario
COPY --from=build /usr/local/lib/libbcg729.so* /usr/local/lib/
COPY --from=build /app/transcoding /usr/local/bin/transcoding

WORKDIR /work
ENTRYPOINT ["/usr/local/bin/transcoding"]

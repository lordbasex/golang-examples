FROM golang:latest as builder

LABEL maintainer="Federico Pereira <fpereira@cnsoluciones.com>"

# Copia el código fuente a la imagen.
COPY . /app

# Establece el directorio de trabajo en la aplicación.
WORKDIR /app

# Compila el binario.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app

# Crea una nueva imagen a partir de scratch (una imagen vacía).
FROM scratch as myapp

LABEL maintainer="Federico Pereira <fpereira@cnsoluciones.com>"

# Copia el binario de la imagen anterior a la imagen scratch.
COPY --from=builder /app/app /app

# Establece el binario como el comando de entrada de la imagen.
ENTRYPOINT ["/app"]

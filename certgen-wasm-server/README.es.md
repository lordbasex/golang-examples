# certgen-wasm-server

[English](README.md) | [Español](README.es.md)

Genera certificados TLS internos directamente en el navegador usando **Astro + Go WASM + Service Worker**.

## Qué hace este proyecto

El formulario web solicita:

- `domain`
- `ip`
- `org`

Luego el backend WASM genera y descarga un ZIP con:

```text
certs/
  ca.crt
  ca.key
  server.crt
  server.key
  fullchain.crt
  install_ca_cert.sh
```

## Vigencia de certificados

- `ca.crt` (certificado CA): **10 años**
- `server.crt` (certificado de servidor): **5 años**

Estos periodos de vigencia se aplican en el backend WASM durante la generación.

## PBX soportados

- Asterisk
- FreePBX
- IssabelPBX
- VitalPBX

Correspondencia genérica de campos:

- **Certificate** -> `server.crt`
- **Key** -> `server.key`
- **Chain** -> `ca.crt`

## Capturas de la UI

Panel principal:

![Panel principal de CertGen](image01.png)

Modal de instalación:

![Modal de instalación de CA](image02.png)

## Stack tecnológico

- Astro (frontend estático)
- Go (`GOOS=js GOARCH=wasm`) para generar certificados
- [`go-wasm-http-server`](https://github.com/nlepage/go-wasm-http-server) dentro de Service Worker

## Estructura del proyecto (actual)

```text
cmd/wasm-server/main.go
internal/certgen/
internal/validation/
pkg/ziputil/
public/
src/
```

## Requisitos

- Go 1.22+ (tienes una versión superior, también funciona)
- Node.js 20+
- npm

## Build y ejecución

```bash
make all
make serve
```

Abrir:

- [http://127.0.0.1:8000](http://127.0.0.1:8000)
- o [http://localhost:8000](http://localhost:8000)

## Idiomas en la UI

El panel soporta:

- Inglés (por defecto)
- Español
- Portugués

## Notas de seguridad

- `ca.key` es altamente sensible (clave privada de la CA raíz).
- `server.key` también es sensible.
- Nunca publiques ni compartas claves privadas en repositorios públicos.


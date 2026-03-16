# Astro + Go WASM Users CRUD (Demo)

[English](README.md) | [Español](README.es.md)

<p>
  <img src="./image_list.png" width="320" alt="Users list" />
  <img src="./image_view.png" width="320" alt="User details" />
  <img src="./image_edit.png" width="320" alt="Edit user" />
</p>

This repository is a **demo** showing how to run a “local API server” inside the browser using:

- **Astro** (UI)
- **Go WebAssembly** + **Service Worker**
- [`go-wasm-http-server`](https://github.com/nlepage/go-wasm-http-server) (HTTP handlers inside a Service Worker)
- **IndexedDB** (local persistence)

The result is a **Users CRUD** where requests like `GET /api/users` are handled by Go WASM in the Service Worker.

## Requirements

- **Node.js** 20+
- **Go** 1.22+

## Install

```bash
make install
```

## Build everything

```bash
make all
```

This will:

- Build the Astro frontend (`dist/`)
- Build the WASM server (`public/server.wasm`)
- Copy the Go runtime file (`public/wasm_exec.js`)

## Run the demo (static server)

```bash
make serve
```

Then open `http://localhost:8000/`.

## API (served by Go WASM)

- `GET /api/users`
- `GET /api/users/{id}`
- `POST /api/users`
- `PUT /api/users/{id}`
- `DELETE /api/users/{id}`

## Validation & errors

Validation happens in:

- the **frontend** (UX)
- the **WASM backend** (authoritative API)
- the **IndexedDB layer** (unique constraints for email/document)

Example API error:

```json
{
  "ok": false,
  "message": "validation_failed",
  "errors": {
    "firstName": "First name is required",
    "email": "Invalid email"
  }
}
```

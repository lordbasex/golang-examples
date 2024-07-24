# WASM CSV Processor

Este proyecto utiliza WebAssembly (WASM) para leer un archivo CSV y crear un JSON con los primeros 10 registros. La aplicación consiste en un archivo HTML simple, CSS para el estilo, y un módulo WASM escrito en Go.

## Explicación

### HTML Estructura

- `input` para cargar el archivo CSV.
- `div` para mostrar el JSON resultante.

### CSS

- Estilos básicos para centrar el contenido y dar formato al `div` de salida.

### JavaScript

- Cargar `wasm_exec.js` y el archivo WASM.
- Configurar un `event listener` para manejar la carga del archivo CSV y procesarlo usando la función `processCSV` del módulo WASM.
- Mostrar el JSON resultante en el `div` de salida.

### Go Código

- Función `processCSV` que lee el CSV, convierte las cabeceras a minúsculas, obtiene los primeros 10 registros y devuelve un JSON.
- Función `main` que configura el entorno para WebAssembly y expone `processCSV` a JavaScript.

## Cómo probarlo

### Compila tu código Go a WASM:

```sh
go mod init csv-wasm
GOOS=js GOARCH=wasm go build -o main.wasm main.go
cp $(go env GOROOT)/misc/wasm/wasm_exec.js .
```

### http-server
```
npm install http-server -g 
```

### test

Abre tu navegador y navega a `http://localhost:8080` (o el puerto que http.server esté usando).
Carga un archivo CSV usando el input y verifica que el JSON se muestre en el div de salida.

```
http-server -p 8080
```

### DOC Explicación

#### CSS: HTML Estructura:

* input para cargar el archivo CSV.
* div para mostrar el JSON resultante.

#### CSS:

* Estilos básicos para centrar el contenido y dar formato al div de salida.

#### JavaScript:

* Cargar wasm_exec.js y el archivo WASM.
* Configurar un event listener para manejar la carga del archivo CSV y procesarlo usando la función processCSV del módulo WASM.
* Mostrar el JSON resultante en el div de salida.

#### Go Código:

* Función `processCSV` que lee el CSV, convierte las cabeceras a minúsculas, obtiene los primeros 10 registros y devuelve un JSON.
* Función main que configura el entorno para WebAssembly y expone processCSV a JavaScript.

#### Cómo probarlo:
Compila tu código Go a WASM:
```sh
GOOS=js GOARCH=wasm go build -o main.wasm main.go
```

Coloca los archivos wasm_exec.js, main.wasm y index.html en el mismo directorio.

Inicia un servidor web para servir estos archivos. Puedes usar un servidor web simple como http-server:

```sh
http-server -p 8080
```
Abre tu navegador y navega a http://localhost:8080 (o el puerto que http.server esté usando).

Carga un archivo CSV usando el input y verifica que el JSON se muestre en el div de salida.


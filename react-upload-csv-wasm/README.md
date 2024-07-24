# React Upload CSV WASM Processor

Este proyecto utiliza React y WebAssembly (WASM) para leer un archivo CSV y crear un JSON con los primeros 10 registros. La aplicación está escrita en React para la interfaz de usuario y Go para la lógica de procesamiento del CSV en WebAssembly.

## Estructura del Proyecto

### HTML

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

## Instalación

### Prerrequisitos

Asegúrate de tener instalado [Node.js](https://nodejs.org/) y [Go](https://golang.org/).

### Clonar el Repositorio

```sh
git clone https://github.com/lordbasex/golang-examples.git
cd react-upload-csv-wasm
````

### Instalar Dependencias
```sh
npm install
```

Este comando instalará los siguientes paquetes:

* react-dropzone: Para manejar la subida de archivos mediante drag and drop.
* papaparse: Para parsear los archivos CSV.
* react-json-view-lite: Para mostrar el JSON de manera interactiva.
* @webassemblyjs/wasm-parser: Para soporte de WebAssembly.

### Construir el Módulo WASM

```sh
cd csv-wasm GOOS=js GOARCH=wasm go build -o ../public/main.wasm main.go && cd .. && npm start
```

Este comando compilará el archivo Go a WebAssembly y colocará el resultado en la carpeta public del proyecto React.

### Ejecución

Para iniciar el servidor de desarrollo de React, usa el siguiente comando:

```sh
npm start
```

Esto abrirá tu aplicación en el navegador, normalmente en `http://localhost:3000`.

### Uso
- Carga un archivo CSV: Usa el input para seleccionar un archivo CSV.
- Visualiza el JSON: El JSON resultante se mostrará en el div de salida.


### Recursos Adicionales
* Go WebAssembly Documentation
* WebAssembly Documentation

Con esta guía, deberías ser capaz de configurar y ejecutar tu aplicación para procesar archivos CSV usando WebAssembly y Go en React.


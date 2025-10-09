# 🎵 wav2g729 - WAV to G.729 Transcoder v1.0.0

Federico Pereira <fpereira@cnsoluciones.com>

Conversor de audio de formato WAV a G.729 utilizando la librería [bcg729](https://github.com/BelledonneCommunications/bcg729) desde Go mediante CGO.

## 🌐 Versiones de idioma
- 🇺🇸 [English](README.md)
- 🇪🇸 [Español](README.es.md) (Actual)

## 📋 Descripción

Este proyecto proporciona una herramienta de línea de comandos que convierte archivos de audio en formato WAV (PCM) a archivos codificados en G.729. El codec G.729 es ampliamente utilizado en telefonía VoIP por su excelente relación entre calidad de voz y compresión.

### Características

- ✅ Conversión de WAV (mono, 8kHz, 16-bit PCM) a G.729
- ✅ Utiliza la librería C `libbcg729` de alta calidad
- ✅ Implementado en Go con CGO para máximo rendimiento
- ✅ Dockerizado para facilitar el despliegue
- ✅ Imagen Docker multi-stage para tamaño optimizado

## 🔧 Requisitos

### Para usar con Docker (Recomendado)
- Docker instalado en tu sistema
- Archivos WAV con las siguientes especificaciones:
  - **Formato**: PCM (AudioFormat = 1)
  - **Canales**: Mono (1 canal)
  - **Sample Rate**: 8000 Hz
  - **Bits por muestra**: 16-bit

### Para compilación local
- Go 1.23 o superior
- CGO habilitado (`CGO_ENABLED=1`)
- `libbcg729` instalada en el sistema
- Herramientas de compilación (gcc, cmake, git)

## 🚀 Uso rápido con Docker

### 📦 Imagen pública disponible

La imagen está disponible públicamente en Docker Hub como `cnsoluciones/wav2g729:latest`. No necesitas construir la imagen localmente.

### 1. Usar la imagen pública (Recomendado)

```bash
docker run --rm -v $PWD:/work cnsoluciones/wav2g729:latest input.wav output.g729
```

**Explicación del comando:**
- `docker run`: Ejecuta un contenedor desde la imagen
- `--rm`: Elimina automáticamente el contenedor después de la ejecución
- `-v $PWD:/work`: Monta el directorio actual en `/work` dentro del contenedor
- `cnsoluciones/wav2g729:latest`: Imagen pública de Docker Hub
- `input.wav`: Archivo de entrada (WAV)
- `output.g729`: Archivo de salida (G.729 raw bitstream)

### 2. Construir la imagen localmente (Opcional)

```bash
docker build -t wav2g729:latest .
```

Este comando:
- Descarga e instala todas las dependencias necesarias
- Compila la librería `bcg729` desde el código fuente
- Compila el programa Go con soporte CGO
- Crea una imagen optimizada de **~19MB** (Alpine Linux)

### 3. Obtener ayuda

```bash
# Mostrar ayuda completa
docker run --rm cnsoluciones/wav2g729:latest --help

# Mostrar versión
docker run --rm cnsoluciones/wav2g729:latest --version

# O simplemente ejecutar sin argumentos
docker run --rm cnsoluciones/wav2g729:latest
```

El helper incluye (en inglés):
- ✅ Descripción del programa y requisitos
- ✅ Ejemplos de uso con Docker
- ✅ Comandos FFmpeg para conversión de archivos incompatibles
- ✅ Información técnica del codec G.729
- ✅ Comandos para verificar la conversión

### 4. Ejemplo con ruta completa

```bash
docker run --rm -v /ruta/a/tus/archivos:/work cnsoluciones/wav2g729:latest audio.wav audio.g729
```

## ✅ Verificar la conversión

Para validar que el archivo `.g729` se creó correctamente, puedes convertirlo de vuelta a WAV con FFmpeg:

```bash
ffmpeg -f g729 -i output.g729 -ar 8000 -ac 1 -c:a pcm_s16le output.wav
```

**Explicación del comando:**
- `-f g729`: Especifica el formato de entrada como G.729 raw bitstream
- `-i output.g729`: Archivo de entrada (el G.729 generado)
- `-ar 8000`: Frecuencia de muestreo de salida (8000 Hz)
- `-ac 1`: Número de canales de salida (1 = mono)
- `-c:a pcm_s16le`: Codec de audio de salida (PCM 16-bit little-endian)
- `output.wav`: Archivo WAV resultante

Ahora puedes reproducir `output.wav` con cualquier reproductor de audio para verificar la calidad de la conversión. Si escuchas el audio correctamente, ¡la conversión fue exitosa! 🎵

## 📁 Estructura del proyecto

```
.
├── Dockerfile           # Definición de la imagen Docker multi-stage
├── go.mod              # Dependencias del proyecto Go
├── go.sum              # Checksums de las dependencias
├── transcoding.go      # Código principal del conversor
└── README.md           # Este archivo
```

## 🏗️ Arquitectura técnica

### Dockerfile Multi-stage optimizado

El proyecto utiliza un Dockerfile de dos etapas optimizado con **Alpine Linux**:

1. **Stage 1 (build)**: Imagen basada en `golang:1.23-alpine`
   - Instala herramientas de compilación (build-base, cmake, git)
   - Clona y compila `bcg729` como librería compartida (`libbcg729.so`)
   - Descarga dependencias Go
   - Compila el binario con CGO habilitado

2. **Stage 2 (runtime)**: Imagen basada en `alpine:latest`
   - Solo contiene el binario compilado y librerías necesarias
   - Copia la librería `libbcg729.so` y dependencias mínimas
   - **Resultado: imagen ultra-ligera de ~19MB** 🚀

### 🎯 Optimizaciones implementadas:

- ✅ **Alpine Linux**: Base mínima (~3MB) vs Debian (~80MB)
- ✅ **Librería compartida**: `libbcg729.so` en lugar de estática
- ✅ **Dependencias mínimas**: Solo `ca-certificates` y `libc6-compat`
- ✅ **Multi-stage build**: Compilación separada del runtime
- ✅ **Sin herramientas de desarrollo**: Solo lo necesario para ejecutar

### Código Go con CGO

El programa utiliza CGO para llamar funciones C de `libbcg729`:

```go
/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lbcg729 -Wl,-rpath,/usr/local/lib
#include <bcg729/encoder.h>
*/
import "C"
```

**Proceso de conversión:**
1. Lee el archivo WAV usando `github.com/youpy/go-wav`
2. Valida el formato (mono, 8kHz, 16-bit PCM)
3. Procesa el audio en frames de 80 muestras (10ms @ 8kHz)
4. Codifica cada frame con `bcg729Encoder`
5. Escribe el bitstream G.729 al archivo de salida

### 🆘 Sistema de ayuda integrado

El programa incluye un helper completo que se activa cuando:
- Se ejecuta sin argumentos: `docker run --rm cnsoluciones/wav2g729:latest`
- Se solicita ayuda explícita: `docker run --rm cnsoluciones/wav2g729:latest --help`

**Características del helper (en inglés):**
- 📋 **Descripción completa** del programa y su propósito
- 📝 **Requisitos técnicos** del archivo WAV de entrada
- 💡 **Ejemplos prácticos** de uso con Docker
- 🔧 **Comandos FFmpeg** para convertir archivos incompatibles
- ✅ **Comandos de verificación** para validar la conversión
- 📊 **Información técnica** del codec G.729
- 🔗 **Enlaces a documentación** adicional

## 🔍 Validación del formato WAV

El programa valida automáticamente que el archivo WAV cumpla con los requisitos:

```
✅ AudioFormat = 1 (PCM)
✅ NumChannels = 1 (Mono)
✅ SampleRate = 8000 Hz
✅ BitsPerSample = 16
```

Si tu archivo no cumple estos requisitos, puedes convertirlo con FFmpeg:

```bash
# Convertir cualquier archivo de audio a formato compatible
ffmpeg -i entrada.mp3 -ar 8000 -ac 1 -sample_fmt s16 salida.wav
```

## 🛠️ Compilación local (sin Docker)

Si prefieres compilar localmente sin Docker:

### 1. Instalar bcg729

```bash
git clone https://github.com/BelledonneCommunications/bcg729
cd bcg729
cmake -S . -B build
cmake --build build --target install
sudo ldconfig
```

### 2. Compilar el programa

```bash
export CGO_ENABLED=1
go build -o transcoding transcoding.go
```

### 3. Ejecutar

```bash
./transcoding input.wav output.g729
```

## 📊 Detalles técnicos del codec G.729

- **Bitrate**: ~8 kbps (muy eficiente)
- **Frame size**: 10ms (80 muestras @ 8kHz)
- **Frame encoding**: ~10 bytes por frame de voz
- **Uso**: VoIP, telefonía IP, videoconferencia
- **Ventaja**: Excelente calidad de voz con mínimo ancho de banda

### VAD (Voice Activity Detection)

El encoder está configurado con VAD deshabilitado (`enableVAD = 0`):
- **VAD = 0**: Todos los frames se codifican como voz (más simple)
- **VAD = 1**: Detecta silencios y los codifica eficientemente (ahorra bandwidth)

Puedes modificar esta configuración en `transcoding.go` línea 19.

## 🐛 Solución de problemas

### Error: "WAV no PCM"
Tu archivo está en formato comprimido. Conviértelo con FFmpeg:
```bash
ffmpeg -i archivo.wav -acodec pcm_s16le salida.wav
```

### Error: "se requiere mono (1 canal)"
Tu archivo es estéreo. Conviértelo a mono:
```bash
ffmpeg -i archivo.wav -ac 1 salida.wav
```

### Error: "se requiere 8000 Hz"
Cambia la frecuencia de muestreo:
```bash
ffmpeg -i archivo.wav -ar 8000 salida.wav
```

### Error: "se requiere 16-bit PCM"
Ajusta el formato de muestra:
```bash
ffmpeg -i archivo.wav -sample_fmt s16 salida.wav
```

### Conversión todo-en-uno con FFmpeg
```bash
ffmpeg -i entrada.mp3 -ar 8000 -ac 1 -sample_fmt s16 -acodec pcm_s16le salida.wav
```

## 📝 Notas importantes

- ⚠️ El archivo de salida `.g729` es un **raw bitstream** sin contenedor
- ⚠️ Para reproducir archivos G.729, necesitas un reproductor compatible o convertirlos de vuelta a WAV
- 💡 **Tip**: Usa `ffmpeg -f g729 -i output.g729 -ar 8000 -ac 1 -c:a pcm_s16le output.wav` para convertir G.729 a WAV
- ⚠️ Algunos codecs G.729 están sujetos a patentes (verifica en tu jurisdicción)
- ⚠️ `bcg729` es una implementación de código abierto y libre de regalías

## 📚 Referencias

- [bcg729 - Biblioteca codec G.729](https://github.com/BelledonneCommunications/bcg729)
- [go-wav - Parser WAV para Go](https://github.com/youpy/go-wav)
- [ITU-T G.729 Specification](https://www.itu.int/rec/T-REC-G.729)
- [CGO Documentation](https://pkg.go.dev/cmd/cgo)

## 📄 Licencia

Este proyecto utiliza `bcg729` que está bajo licencia BSD-like. Verifica los términos de licencia antes de usar en producción.

## 🤝 Contribuciones

Las contribuciones son bienvenidas. Por favor:
1. Haz fork del repositorio
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 👨‍💻 Autor

**Federico Pereira** <fpereira@cnsoluciones.com>

Proyecto de conversión de audio WAV a G.729 usando Go y CGO.

### 🏢 CNSoluciones

Este proyecto es parte de CNSoluciones, especializada en soluciones de telecomunicaciones y VoIP.

---

**¿Preguntas o problemas?** Abre un issue en el repositorio.

## 🌐 Versiones de idioma

- 🇺🇸 [English](README.md)
- 🇪🇸 [Español](README.es.md) (Actual)


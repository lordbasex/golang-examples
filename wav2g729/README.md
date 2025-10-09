# 🎵 wav2g729 - WAV to G.729 Transcoder v1.0.0

Federico Pereira <fpereira@cnsoluciones.com>

Audio transcoder from WAV format to G.729 using the [bcg729](https://github.com/BelledonneCommunications/bcg729) library from Go via CGO.

## 📋 Description

This project provides a command-line tool that converts audio files from WAV (PCM) format to G.729 encoded files. The G.729 codec is widely used in VoIP telephony for its excellent voice quality to compression ratio.

### Features

- ✅ WAV conversion (mono, 8kHz, 16-bit PCM) to G.729
- ✅ Uses high-quality C `libbcg729` library
- ✅ Implemented in Go with CGO for maximum performance
- ✅ Dockerized for easy deployment
- ✅ Multi-stage Docker image for optimized size

## 🔧 Requirements

### For Docker usage (Recommended)
- Docker installed on your system
- WAV files with the following specifications:
  - **Format**: PCM (AudioFormat = 1)
  - **Channels**: Mono (1 channel)
  - **Sample Rate**: 8000 Hz
  - **Bits per sample**: 16-bit

### For local compilation
- Go 1.23 or higher
- CGO enabled (`CGO_ENABLED=1`)
- `libbcg729` installed on the system
- Build tools (gcc, cmake, git)

## 🚀 Quick usage with Docker

### 📦 Public image available

The image is publicly available on Docker Hub as `cnsoluciones/wav2g729:latest`. You don't need to build the image locally.

### 1. Use the public image (Recommended)

```bash
docker run --rm -v $PWD:/work cnsoluciones/wav2g729:latest input.wav output.g729
```

**Command explanation:**
- `docker run`: Runs a container from the image
- `--rm`: Automatically removes the container after execution
- `-v $PWD:/work`: Mounts the current directory to `/work` inside the container
- `cnsoluciones/wav2g729:latest`: Public Docker Hub image
- `input.wav`: Input file (WAV)
- `output.g729`: Output file (G.729 raw bitstream)

### 2. Build the image locally (Optional)

```bash
docker build -t wav2g729:latest .
```

This command:
- Downloads and installs all necessary dependencies
- Compiles the `bcg729` library from source code
- Compiles the Go program with CGO support
- Creates an optimized image of **~19MB** (Alpine Linux)

### 3. Get help

```bash
# Show complete help
docker run --rm cnsoluciones/wav2g729:latest --help

# Show version
docker run --rm cnsoluciones/wav2g729:latest --version

# Or simply run without arguments
docker run --rm cnsoluciones/wav2g729:latest
```

The helper includes (in English):
- ✅ Program description and requirements
- ✅ Docker usage examples
- ✅ FFmpeg commands for incompatible file conversion
- ✅ Technical information about G.729 codec
- ✅ Commands to verify conversion

### 4. Example with full path

```bash
docker run --rm -v /path/to/your/files:/work cnsoluciones/wav2g729:latest audio.wav audio.g729
```

## ✅ Verify conversion

To validate that the `.g729` file was created correctly, you can convert it back to WAV with FFmpeg:

```bash
ffmpeg -f g729 -i output.g729 -ar 8000 -ac 1 -c:a pcm_s16le output.wav
```

**Command explanation:**
- `-f g729`: Specifies input format as G.729 raw bitstream
- `-i output.g729`: Input file (the generated G.729)
- `-ar 8000`: Output sample rate (8000 Hz)
- `-ac 1`: Number of output channels (1 = mono)
- `-c:a pcm_s16le`: Output audio codec (PCM 16-bit little-endian)
- `output.wav`: Resulting WAV file

Now you can play `output.wav` with any audio player to verify the conversion quality. If you hear the audio correctly, the conversion was successful! 🎵

## 📁 Project structure

```
.
├── Dockerfile           # Multi-stage Docker image definition
├── go.mod              # Go project dependencies
├── go.sum              # Dependency checksums
├── transcoding.go      # Main transcoder code
├── README.md           # This file (English)
└── README.es.md        # Spanish documentation
```

## 🏗️ Technical architecture

### Optimized Multi-stage Dockerfile

The project uses a two-stage Dockerfile optimized with **Alpine Linux**:

1. **Stage 1 (build)**: Image based on `golang:1.23-alpine`
   - Installs build tools (build-base, cmake, git)
   - Clones and compiles `bcg729` as shared library (`libbcg729.so`)
   - Downloads Go dependencies
   - Compiles binary with CGO enabled

2. **Stage 2 (runtime)**: Image based on `alpine:latest`
   - Contains only the compiled binary and necessary libraries
   - Copies `libbcg729.so` and minimal dependencies
   - **Result: ultra-light image of ~19MB** 🚀

### 🎯 Implemented optimizations:

- ✅ **Alpine Linux**: Minimal base (~3MB) vs Debian (~80MB)
- ✅ **Shared library**: `libbcg729.so` instead of static
- ✅ **Minimal dependencies**: Only `ca-certificates` and `libc6-compat`
- ✅ **Multi-stage build**: Separate compilation from runtime
- ✅ **No development tools**: Only what's necessary to run

### Go code with CGO

The program uses CGO to call C functions from `libbcg729`:

```go
/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lbcg729 -Wl,-rpath,/usr/local/lib
#include <bcg729/encoder.h>
*/
import "C"
```

**Conversion process:**
1. Reads WAV file using `github.com/youpy/go-wav`
2. Validates format (mono, 8kHz, 16-bit PCM)
3. Processes audio in frames of 80 samples (10ms @ 8kHz)
4. Encodes each frame with `bcg729Encoder`
5. Writes G.729 bitstream to output file

### 🆘 Integrated help system

The program includes a complete helper that activates when:
- Run without arguments: `docker run --rm cnsoluciones/wav2g729:latest`
- Explicit help request: `docker run --rm cnsoluciones/wav2g729:latest --help`

**Helper features (in English):**
- 📋 **Complete description** of the program and its purpose
- 📝 **Technical requirements** for input WAV file
- 💡 **Practical examples** of Docker usage
- 🔧 **FFmpeg commands** to convert incompatible files
- ✅ **Verification commands** to validate conversion
- 📊 **Technical information** about G.729 codec
- 🔗 **Additional documentation** links

## 🔍 WAV format validation

The program automatically validates that the WAV file meets the requirements:

```
✅ AudioFormat = 1 (PCM)
✅ NumChannels = 1 (Mono)
✅ SampleRate = 8000 Hz
✅ BitsPerSample = 16
```

If your file doesn't meet these requirements, you can convert it with FFmpeg:

```bash
# Convert any audio file to compatible format
ffmpeg -i input.mp3 -ar 8000 -ac 1 -sample_fmt s16 -acodec pcm_s16le output.wav
```

## 🛠️ Local compilation (without Docker)

If you prefer to compile locally without Docker:

### 1. Install bcg729

```bash
git clone https://github.com/BelledonneCommunications/bcg729
cd bcg729
cmake -S . -B build
cmake --build build --target install
sudo ldconfig
```

### 2. Compile the program

```bash
export CGO_ENABLED=1
go build -o transcoding transcoding.go
```

### 3. Run

```bash
./transcoding input.wav output.g729
```

## 📊 G.729 codec technical details

- **Bitrate**: ~8 kbps (very efficient)
- **Frame size**: 10ms (80 samples @ 8kHz)
- **Frame encoding**: ~10 bytes per voice frame
- **Usage**: VoIP, IP telephony, videoconferencing
- **Advantage**: Excellent voice quality with minimal bandwidth

### VAD (Voice Activity Detection)

The encoder is configured with VAD disabled (`enableVAD = 0`):
- **VAD = 0**: All frames are encoded as voice (simpler)
- **VAD = 1**: Detects silence and encodes it efficiently (saves bandwidth)

You can modify this configuration in `transcoding.go` line 19.

## 🐛 Troubleshooting

### Error: "WAV is not PCM"
Your file is in compressed format. Convert it with FFmpeg:
```bash
ffmpeg -i file.wav -acodec pcm_s16le output.wav
```

### Error: "mono required (1 channel)"
Your file is stereo. Convert to mono:
```bash
ffmpeg -i file.wav -ac 1 output.wav
```

### Error: "8000 Hz required"
Change the sample rate:
```bash
ffmpeg -i file.wav -ar 8000 output.wav
```

### Error: "16-bit PCM required"
Adjust the sample format:
```bash
ffmpeg -i file.wav -sample_fmt s16 output.wav
```

### All-in-one conversion with FFmpeg
```bash
ffmpeg -i input.mp3 -ar 8000 -ac 1 -sample_fmt s16 -acodec pcm_s16le output.wav
```

## 📝 Important notes

- ⚠️ The output `.g729` file is a **raw bitstream** without container
- ⚠️ To play G.729 files, you need a compatible player or convert them back to WAV
- 💡 **Tip**: Use `ffmpeg -f g729 -i output.g729 -ar 8000 -ac 1 -c:a pcm_s16le output.wav` to convert G.729 to WAV
- ⚠️ Some G.729 codecs are subject to patents (check in your jurisdiction)
- ⚠️ `bcg729` is an open-source and royalty-free implementation

## 📚 References

- [bcg729 - G.729 codec library](https://github.com/BelledonneCommunications/bcg729)
- [go-wav - WAV parser for Go](https://github.com/youpy/go-wav)
- [ITU-T G.729 Specification](https://www.itu.int/rec/T-REC-G.729)
- [CGO Documentation](https://pkg.go.dev/cmd/cgo)

## 📄 License

This project uses `bcg729` which is under BSD-like license. Check license terms before using in production.

## 🤝 Contributing

Contributions are welcome. Please:
1. Fork the repository
2. Create a branch for your feature (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 👨‍💻 Author

**Federico Pereira** <fpereira@cnsoluciones.com>

WAV to G.729 audio conversion project using Go and CGO.

### 🏢 CNSoluciones

This project is part of CNSoluciones, specialized in telecommunications and VoIP solutions.

---

**Questions or issues?** Open an issue in the repository.

## 🌐 Language versions

- 🇺🇸 [English](README.md) (Current)
- 🇪🇸 [Español](README.es.md)

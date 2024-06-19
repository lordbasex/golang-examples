# WASM MP3 TO WAV

## BUILD
```
GOOS=js GOARCH=wasm go build -o main.wasm main.go
```

## RUN 
```
python3 -m http.server 8080
```
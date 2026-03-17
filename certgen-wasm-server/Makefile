WASM=public/server.wasm
WASM_EXEC=public/wasm_exec.js
PORT?=8000

.PHONY: install build-frontend build-wasm runtime all dev preview serve clean

install:
	npm install

node_modules:
	npm install

build-frontend: node_modules
	npm run build

build-wasm:
	GOOS=js GOARCH=wasm go build -o $(WASM) ./cmd/wasm-server

runtime:
	cp $$(go env GOROOT)/lib/wasm/wasm_exec.js $(WASM_EXEC) || cp $$(go env GOROOT)/misc/wasm/wasm_exec.js $(WASM_EXEC)

all: build-wasm runtime build-frontend

dev:
	npm run dev

preview:
	npm run preview

serve:
	python3 -m http.server $(PORT) -d dist --bind 127.0.0.1

clean:
	rm -rf dist
	rm -f $(WASM) $(WASM_EXEC)


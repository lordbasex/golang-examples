# GO WAV-CHECK WASM

## Build Instructions

1. Run the following command to build the WebAssembly (WASM) binary:

    ```bash
    make
    ```

2. Install an HTTP server if not already installed. You can use `http-server`, which can be installed globally using npm:

    ```bash
    npm install --global http-server
    ```

    or on macOS using Homebrew:

    ```bash
    brew install http-server
    ```

3. Start the HTTP server:

    ```bash
    http-server
    ```

4. Open your web browser and navigate to [http://127.0.0.1:8080/](http://127.0.0.1:8080/) to access the site.

## About

This project demonstrates using WebAssembly to check properties of WAV files. The WebAssembly code is written in Go and compiled to WASM.

## Credits

- [Go Audio](https://github.com/go-audio) - Library used for handling audio files in Go.
- [GitHub Pages](https://github.com/lordbasex/golang-examples/) - Hosting service for the project site.

Feel free to contribute or report issues!

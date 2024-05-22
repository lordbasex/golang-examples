# GO WEBSSH TO WEBSOCKET


This project allows you to connect to an SSH server over WebSocket directly from your browser.

It is based on the GO SSH TO WEBSOCKET project by Razikus but has been redesigned to use the Fiber framework and additional libraries.





## Features

- SSH connection via WebSockets using the [Fiber framework](https://docs.gofiber.io/)
- Web-based terminal emulation using [xterm.js](https://xtermjs.org/)
- Environment variable configuration for SSH credentials and settings.
- Dockerized for easy deployment.


## Libraries Used

### Golang

- Fiber: A web framework inspired by Express.js.
- Gorilla WebSocket: A library for handling WebSockets in Go.
- SSH: A Go package for handling SSH connections.

### JavaScript

- xterm.js: Embedded terminal for web applications.
- xterm-addon-fit: An addon for xterm.js to fit the terminal size.


## Running in docker

To run the project using Docker:

```
docker run -p 8280:8280 -e SSH_USER=USER -e SSH_PASS=PASS -e SSH_HOST=HOST -e SSH_PORT=PORT --rm lordbasex/webssh:latest
```

Then, go to [http://localhost:8280](http://localhost:8280) and you will see the terminal in your browser.


## Configuration

Available environment variables:

```
SSH_USER="your_username"
SSH_PASS="your_password"
SSH_HOST="ssh.example.com"
SSH_PORT="22"
```

We hope this project is useful to you. If you have any questions or suggestions, feel free to contact us.
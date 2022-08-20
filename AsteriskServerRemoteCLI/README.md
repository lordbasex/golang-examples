# Asterisk Server Remote CLI


## RABBITMQ SERVER - Docker-Compose

```docker-compose.yml
version: "3.7"
services:
  rabbitmq:
    build:
      context: ./
      dockerfile: Dockerfile
    container_name: rabbitmq
    restart: always
    ports:
      - 7777:15672
      - 8888:5672
    hostname: rabbitmq
    volumes:
      - ./data:/var/lib/rabbitmq/mnesia
    environment:
      - TZ=America/Argentina/Buenos_Aires
      - RABBITMQ_DEFAULT_USER=rabbitUser
      - RABBITMQ_DEFAULT_PASS=AASwPslfkjJs2ijsnfiujhaADXKjbsadkjbdasdc222asd11A
    networks:
      rabbitmq_net:
        aliases:
          - rabbitmq_host

volumes:
  data: {}

networks:
  rabbitmq_net:
    name: rabbitmq_network
    driver: bridge
```

## RUN DEV AsteriskServerRemoteCLI
```
go mod init AsteriskServerRemoteCLI
go mod tidy
go run main.go
```

## BUILD SERVER AsteriskServerRemoteCLI
```
go build -o asterisk-server-remote-cli main.go

```

## systemctl - Debian and Ubuntu
```
yes|cp -fra asterisk-server-remote-cli /usr/local/bin/
```

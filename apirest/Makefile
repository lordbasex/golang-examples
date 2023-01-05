# Nombre de la imagen de Docker.
IMAGE_NAME = cnsoluciones/myapp

# Crea el binario de la aplicación.
build:
	go build

# Crea la imagen de Docker de la aplicación.
docker-build: build
	docker build -t $(IMAGE_NAME) .

# Ejecuta la reseta de docker-compose.yml haciendo un deploy de myapp y mysql.
docker-run:
	docker-compose up -d 

# Nombre de la imagen de Docker.
IMAGE_NAME = cnsoluciones/myapp

# Crea la imagen de Docker de la aplicación.
docker-build:
	docker build -t $(IMAGE_NAME) .

# Ejecuta la reseta de docker-compose.yml haciendo un deploy de myapp y mysql.
docker-run:
	docker-compose up -d 

# Utilizamos la imagen base de Golang con Alpine
FROM golang:1.22.3-alpine3.20

# Instalar bash para usar scripts en el contenedor
RUN apk add --no-cache bash

# Establecemos el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiamos el go.mod y go.sum primero y descargamos los módulos necesarios
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto de los archivos fuente y el script wait-for-it.sh
COPY . .
COPY wait-for-it.sh /wait-for-it.sh

# Hacemos que el script sea ejecutable y construimos la aplicación
RUN chmod +x /wait-for-it.sh && \
    go build -o /app/bin/myapp ./cmd/rest/main.go

# Exponemos el puerto 8080 para la aplicación
EXPOSE 8080

# Comando para ejecutar la aplicación
# Usamos wait-for-it.sh para asegurarnos de que MySQL esté disponible antes de iniciar la aplicación
CMD ["/wait-for-it.sh", "mysql:3306", "--", "/app/bin/myapp"]

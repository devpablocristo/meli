#!/bin/sh

# Cargar las variables de entorno desde el archivo .env
if [ -f .env ]; then
  source .env
fi

# Imprimir las variables de entorno para verificar
echo "MYSQL_CONT: ${MYSQL_CONT}"
echo "MYSQL_CONT_PORT: ${MYSQL_CONT_PORT}"

# Utilizar las variables de entorno en el comando wait-for
wait-for "${MYSQL_CONT}:${MYSQL_CONT_PORT}" -- "$@"

# Otros comandos aquí (si es necesario)
CompileDaemon --build="go build -o ./build/tarefaapi -v ./cmd/api" -command="./build/tarefaapi"

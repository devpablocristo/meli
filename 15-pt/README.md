# Manual de Usuario para la API de Inventario

Este manual proporciona instrucciones detalladas sobre cómo configurar y ejecutar la API de Inventario, desarrollada en Golang utilizando MySQL como base de datos y gestionada mediante Docker. Asegúrese de seguir cada paso cuidadosamente para garantizar una correcta configuración y funcionamiento de la aplicación.

## Prerrequisitos

1. Docker instalado en su sistema.
2. Docker Compose instalado en su sistema.
3. Conexión a internet para descargar las imágenes de Docker necesarias.

## Contenido del Repositorio

- `Dockerfile`: Archivo de configuración para construir la imagen de la aplicación Golang.
- `docker-compose.yml`: Archivo de configuración para orquestar los servicios Docker (aplicación, MySQL y phpMyAdmin).
- `init.sql`: Script SQL para inicializar la base de datos MySQL con el esquema necesario y el usuario de la API.
- Código fuente de la API de Inventario.

## Instrucciones de Configuración

### Paso 1: Configurar la Base de Datos

Antes de iniciar los servicios Docker, es necesario asegurarse de que el script `init.sql` se ejecute para configurar la base de datos. Este script crea la base de datos `inventory`, la tabla `items` y un usuario de la API con los permisos adecuados.

### Paso 2: Iniciar los Servicios Docker

Utilice Docker Compose para iniciar todos los servicios definidos en el archivo `docker-compose.yml`.

```sh
docker-compose up --build
```

Este comando hará lo siguiente:

1. Construirá la imagen de la aplicación Golang.
2. Iniciará el contenedor de MySQL.
3. Iniciará el contenedor de phpMyAdmin.
4. Iniciará el contenedor de la aplicación Golang.

### Paso 3: Ejecutar el Script SQL en phpMyAdmin

Abra su navegador web y vaya a `http://localhost:8081` para acceder a phpMyAdmin. Inicie sesión con las siguientes credenciales:

- Usuario: `root`
- Contraseña: `root`

Una vez dentro de phpMyAdmin, siga estos pasos:

1. Seleccione la base de datos `inventory`.
2. Vaya a la pestaña "SQL".
3. Copie y pegue el contenido del archivo `init.sql`.
4. Ejecute el script.

Esto inicializará la base de datos y configurará el usuario necesario para la API.

### Paso 4: Verificar el Funcionamiento de la API

Una vez que todos los contenedores estén en funcionamiento y la base de datos esté configurada, puede verificar el funcionamiento de la API accediendo a `http://localhost:8080` en su navegador web o utilizando herramientas como `curl` o `Postman` para interactuar con los endpoints `/items`.

### Endpoints de la API

- **POST /items**: Crear un nuevo ítem en el inventario.
  - Body JSON:
    ```json
    {
      "code": "1234",
      "title": "Sample Item",
      "description": "This is a sample item",
      "price": 19.99,
      "stock": 100,
      "status": "available"
    }
    ```

- **GET /items**: Obtener una lista de todos los ítems en el inventario.

### Ejemplo de Uso con `curl`

#### Crear un Nuevo Ítem

```sh
curl -X POST http://localhost:8080/items -H "Content-Type: application/json" -d '{
  "code": "1234",
  "title": "Sample Item",
  "description": "This is a sample item",
  "price": 19.99,
  "stock": 100,
  "status": "available"
}'
```

#### Obtener la Lista de Ítems

```sh
curl http://localhost:8080/items
```

## Solución de Problemas

### Error de Conexión a MySQL

Si la aplicación Golang no puede conectarse a MySQL, asegúrese de que:

1. El contenedor de MySQL esté en funcionamiento.
2. El script `init.sql` se haya ejecutado correctamente.
3. Las credenciales de la base de datos en el código de configuración coincidan con las del script `init.sql`.

### Verificación de Logs

Puede verificar los logs de los contenedores Docker para obtener más detalles sobre cualquier error.

```sh
docker-compose logs app
docker-compose logs mysql
docker-compose logs phpmyadmin
```

## Conclusión

Siguiendo estos pasos, debería poder configurar y ejecutar correctamente la API de Inventario. Si encuentra algún problema, consulte los logs de los contenedores Docker y asegúrese de que todos los servicios estén configurados correctamente.
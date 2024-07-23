# API Inventory

Este manual proporciona instrucciones detalladas sobre cómo configurar y ejecutar la API de Inventario, desarrollada en Golang utilizando MySQL y MongoDB como bases de datos y gestionada mediante Docker. Asegúrese de seguir cada paso cuidadosamente para garantizar una correcta configuración y funcionamiento de la aplicación.

## Prerrequisitos

1. Docker instalado en su sistema.
2. Docker Compose instalado en su sistema.
3. Conexión a internet para descargar las imágenes de Docker necesarias.

## Contenido del Repositorio

- `Dockerfile`: Archivo de configuración para construir la imagen de la aplicación Golang.
- `docker-compose.yml`: Archivo de configuración para orquestar los servicios Docker (aplicación, MySQL, phpMyAdmin, MongoDB y mongo-express).
- `init.sql`: Script SQL para inicializar la base de datos MySQL con el esquema necesario y el usuario de la API.
- `init-mongo.js`: Script JavaScript para inicializar la base de datos MongoDB con el esquema necesario y el usuario de la API.
- Código fuente de la API de Inventario.

## Instrucciones de Configuración

### Paso 1: Configurar la Base de Datos MySQL

Antes de iniciar los servicios Docker, es necesario asegurarse de que el script `init.sql` se ejecute para configurar la base de datos MySQL. Este script crea la base de datos `inventory`, la tabla `items` y un usuario de la API con los permisos adecuados.

### Paso 2: Iniciar los Servicios Docker

Utilice Docker Compose para iniciar todos los servicios definidos en el archivo `docker-compose.yml`.

```sh
docker-compose up --build
```

Este comando hará lo siguiente:

1. Construirá la imagen de la aplicación Golang.
2. Iniciará el contenedor de MySQL.
3. Iniciará el contenedor de phpMyAdmin.
4. Iniciará el contenedor de MongoDB.
5. Iniciará el contenedor de mongo-express.
6. Iniciará el contenedor de la aplicación Golang.

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

### Paso 4: Crear Base de Datos, Usuario y Colección en Mongo-Express

Abra su navegador web y vaya a `http://localhost:8082` para acceder a Mongo-Express. Inicie sesión con las siguientes credenciales:

- Usuario: `root`
- Contraseña: `root`

Una vez dentro de Mongo-Express, siga estos pasos:

1. **Crear la Base de Datos `inventory`**:
   - En la barra lateral izquierda, haz clic en **"Add Database"**.
   - Ingresa `inventory` como el nombre de la nueva base de datos.
   - Haz clic en **"Submit"** para crear la base de datos.

2. **Crear un Usuario para la Base de Datos `inventory`**:
   - Selecciona la base de datos `inventory` en la barra lateral izquierda.
   - Haz clic en la pestaña **"Add User"**.
   - Ingresa el nombre de usuario (`api_user`), la contraseña (`api_password`), y selecciona el rol `readWrite` para permisos de lectura y escritura.
   - Haz clic en **"Submit"** para crear el usuario.

3. **Crear la Colección `items`**:
   - Aún dentro de la base de datos `inventory`, haz clic en **"Add Collection"**.
   - Ingresa `items` como el nombre de la colección.
   - Haz clic en **"Submit"** para crear la colección.

4. **Insertar Datos Iniciales en la Colección `items`**:
   - Selecciona la colección `items` dentro de la base de datos `inventory`.
   - Haz clic en **"Add Document"**.
   - Copia y pega el siguiente contenido JSON en el campo de entrada de documentos:
     ```json
     {
       "id": 1,
       "code": "ITEM001",
       "title": "Example Item",
       "description": "This is an example item",
       "price": 29.99,
       "stock": 50,
       "status": "available",
       "created_at": "2024-07-17T15:04:05Z",
       "updated_at": "2024-07-17T15:04:05Z"
     }
     ```
   - Haz clic en **"Submit"** para insertar el documento.

Esto inicializará la base de datos, configurará el usuario necesario y creará la colección `items` con un documento de ejemplo para la API.

### Paso 5: Verificar el Funcionamiento de la API

Una vez que todos los contenedores estén en funcionamiento y la base de datos esté configurada, puede verificar el funcionamiento de la API accediendo a `http://localhost:8080` en su navegador web o utilizando herramientas como `curl` o `Postman` para interactuar con los endpoints `/items`.

### Endpoints de la API

- **POST /items**: Crear un nuevo ítem en el inventario.
  - Body JSON:
    ```json
    {
      "id": 1,
      "code": "ITEM001",
      "title": "Example Item",
      "description": "This is an example item",
      "price": 29.99,
      "stock": 50,
      "status": "available",
      "created_at": "2024-07-17T15:04:05Z",
      "updated_at": "2024-07-17T15:04:05Z"
    }
    ```

- **GET /items**: Obtener una lista de todos los ítems en el inventario.

### Ejemplo de Uso con `curl`

#### Crear un Nuevo Ítem

```sh
curl -X POST http://localhost:8080/items -H "Content-Type: application/json" -d '{
  "id": 1,
  "code": "ITEM001",
  "title": "Example Item",
  "description": "This is an example item",
  "price": 29.99,
  "stock": 50,
  "status": "available",
  "created_at": "2024-07-17T15:04:05Z",
  "updated_at": "2024-07-17T15:04:05Z"
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
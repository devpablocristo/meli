### Resumen:

Este archivo `docker-compose.yml` configura y ejecuta una aplicación compuesta por tres servicios: una aplicación Go (`api`), una base de datos MySQL (`mysql`) y una herramienta de administración web para MySQL (`phpmyadmin`). Define cómo se deben construir e iniciar estos servicios, cómo se deben comunicar entre ellos y cómo se deben mapear los puertos y volúmenes para su correcto funcionamiento.

### ¿Qué es Docker Compose?

Docker Compose es una herramienta para definir y ejecutar aplicaciones Docker multicontenedor. Utiliza un archivo YAML para configurar los servicios de la aplicación. Con un solo comando (`docker-compose up`), puedes crear e iniciar todos los servicios definidos en el archivo `docker-compose.yml`.

### Explicación del archivo `docker-compose.yml`

Este archivo define una configuración para una aplicación compuesta por tres servicios: `api`, `mysql` y `phpmyadmin`. A continuación, se explica cada sección del archivo:

```yaml
version: "3.8"
```
- **version**: Especifica la versión del formato de archivo de Docker Compose que se está utilizando.

#### Servicios:

```yaml
services:
```
- **services**: Define los contenedores que se iniciarán. Cada servicio representa un contenedor.

##### Servicio `api`:

```yaml
  api:
    container_name: myapp
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    volumes:
      - type: bind
        source: ./
        target: /app
    environment:
      - APP_NAME=myapp
      - DEBUG=true
      - DEV_DB_HOST=mysql
      - DEV_DB_HOST_PORT=3306
      - DEV_DB_DATABASE=dev_events_db
      - DEV_DB_USERNAME=root
      - DEV_DB_USER_PASSWORD=root
    depends_on:
      - mysql
    networks:
      - default
```
- **container_name**: Asigna un nombre al contenedor.
- **build**: Especifica cómo construir la imagen Docker para el servicio `api`.
  - **context**: El directorio de contexto de construcción, en este caso el directorio actual.
  - **dockerfile**: El Dockerfile que se utilizará para construir la imagen.
- **ports**: Mapea el puerto 8080 del contenedor al puerto 8080 del host.
- **volumes**: Monta el directorio actual (`./`) en el contenedor en `/app`.
- **environment**: Define las variables de entorno para el contenedor.
- **depends_on**: Especifica que el servicio `api` depende de que el servicio `mysql` se inicie primero.
- **networks**: Conecta el servicio a la red `default`.

##### Servicio `mysql`:

```yaml
  mysql:
    image: mysql:8.0
    container_name: mysql
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: dev_events_db
      MYSQL_USER: root
      MYSQL_PASSWORD: root
    volumes:
      - mysql_data:/var/lib/mysql
    networks:
      - default
```
- **image**: Utiliza la imagen oficial de MySQL 8.0.
- **container_name**: Asigna un nombre al contenedor.
- **ports**: Mapea el puerto 3306 del contenedor al puerto 3306 del host.
- **environment**: Define las variables de entorno para configurar MySQL.
- **volumes**: Utiliza un volumen llamado `mysql_data` para persistir los datos de MySQL.
- **networks**: Conecta el servicio a la red `default`.

##### Servicio `phpmyadmin`:

```yaml
  phpmyadmin:
    image: phpmyadmin/phpmyadmin:latest
    container_name: phpmyadmin
    ports:
      - "8081:80"
    environment:
      PMA_HOST: mysql
      MYSQL_ROOT_PASSWORD: root
    depends_on:
      - mysql
    networks:
      - default
```
- **image**: Utiliza la imagen oficial de phpMyAdmin.
- **container_name**: Asigna un nombre al contenedor.
- **ports**: Mapea el puerto 80 del contenedor al puerto 8081 del host.
- **environment**: Define las variables de entorno para configurar phpMyAdmin.
- **depends_on**: Especifica que el servicio `phpmyadmin` depende de que el servicio `mysql` se inicie primero.
- **networks**: Conecta el servicio a la red `default`.

#### Redes:

```yaml
networks:
  default:
    driver: bridge
```
- **networks**: Define las redes para los servicios.
  - **default**: Crea una red llamada `default` utilizando el driver `bridge`.

#### Volúmenes:

```yaml
volumes:
  mysql_data:
```
- **volumes**: Define los volúmenes para persistir los datos.
  - **mysql_data**: Crea un volumen llamado `mysql_data` para almacenar los datos de MySQL de manera persistente.


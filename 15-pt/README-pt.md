### API de Inventário

Este manual fornece instruções detalhadas sobre como configurar e executar a API de Inventário, desenvolvida em Golang utilizando MySQL como banco de dados e gerenciada por meio do Docker. Certifique-se de seguir cada passo cuidadosamente para garantir uma configuração e funcionamento corretos da aplicação.

## Pré-requisitos

1. Docker instalado no seu sistema.
2. Docker Compose instalado no seu sistema.
3. Conexão à internet para baixar as imagens necessárias do Docker.

## Conteúdo do Repositório

- `Dockerfile`: Arquivo de configuração para construir a imagem da aplicação Golang.
- `docker-compose.yml`: Arquivo de configuração para orquestrar os serviços Docker (aplicação, MySQL e phpMyAdmin).
- `init.sql`: Script SQL para inicializar o banco de dados MySQL com o esquema necessário e o usuário da API.
- Código-fonte da API de Inventário.

## Instruções de Configuração

### Passo 1: Configurar o Banco de Dados

Antes de iniciar os serviços Docker, é necessário garantir que o script `init.sql` seja executado para configurar o banco de dados. Este script cria o banco de dados `inventory`, a tabela `items` e um usuário da API com as permissões adequadas.

### Passo 2: Iniciar os Serviços Docker

Utilize o Docker Compose para iniciar todos os serviços definidos no arquivo `docker-compose.yml`.

```sh
docker-compose up --build
```

Este comando fará o seguinte:

1. Construirá a imagem da aplicação Golang.
2. Iniciará o contêiner do MySQL.
3. Iniciará o contêiner do phpMyAdmin.
4. Iniciará o contêiner da aplicação Golang.

### Passo 3: Executar o Script SQL no phpMyAdmin

Abra o seu navegador web e vá para `http://localhost:8081` para acessar o phpMyAdmin. Faça login com as seguintes credenciais:

- Usuário: `root`
- Senha: `root`

Uma vez dentro do phpMyAdmin, siga estes passos:

1. Selecione o banco de dados `inventory`.
2. Vá para a aba "SQL".
3. Copie e cole o conteúdo do arquivo `init.sql`.
4. Execute o script.

Isso inicializará o banco de dados e configurará o usuário necessário para a API.

### Passo 4: Verificar o Funcionamento da API

Uma vez que todos os contêineres estejam em funcionamento e o banco de dados esteja configurado, você pode verificar o funcionamento da API acessando `http://localhost:8080` no seu navegador web ou utilizando ferramentas como `curl` ou `Postman` para interagir com os endpoints `/items`.

### Endpoints da API

- **POST /items**: Criar um novo item no inventário.
  - Corpo JSON:
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

- **GET /items**: Obter uma lista de todos os itens no inventário.

### Exemplo de Uso com `curl`

#### Criar um Novo Item

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

#### Obter a Lista de Itens

```sh
curl http://localhost:8080/items
```

## Solução de Problemas

### Erro de Conexão ao MySQL

Se a aplicação Golang não conseguir se conectar ao MySQL, certifique-se de que:

1. O contêiner do MySQL está em funcionamento.
2. O script `init.sql` foi executado corretamente.
3. As credenciais do banco de dados no código de configuração coincidem com as do script `init.sql`.

### Verificação de Logs

Você pode verificar os logs dos contêineres Docker para obter mais detalhes sobre qualquer erro.

```sh
docker-compose logs app
docker-compose logs mysql
docker-compose logs phpmyadmin
```

## Conclusão

Seguindo estes passos, você deverá ser capaz de configurar e executar corretamente a API de Inventário. Se encontrar algum problema, consulte os logs dos contêineres Docker e certifique-se de que todos os serviços estão configurados corretamente.
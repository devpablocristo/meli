# Utilizamos a imagem base do Golang com Alpine
FROM golang:1.22.3-alpine3.20

# Instalar bash para usar scripts no contêiner
RUN apk add --no-cache bash

# Estabelecemos o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copiamos o go.mod e go.sum primeiro e baixamos os módulos necessários
COPY go.mod go.sum ./
RUN go mod download

# Copiamos o restante dos arquivos fonte e o script wait-for-it.sh
COPY . .
COPY wait-for-it.sh /wait-for-it.sh

# Tornamos o script executável e construímos a aplicação
RUN chmod +x /wait-for-it.sh && \
    go build -o /app/bin/myapp ./cmd/rest/main.go

# Exponhamos a porta 8080 para a aplicação
EXPOSE 8080

# Comando para executar a aplicação
# Usamos wait-for-it.sh para garantir que o MySQL esteja disponível antes de iniciar a aplicação
CMD ["/wait-for-it.sh", "mysql:3306", "--", "/app/bin/myapp"]
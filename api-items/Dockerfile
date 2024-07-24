# Utilizamos a imagem base do Golang com Alpine
FROM golang:1.22.3-alpine3.20

# Instalar bash para usar scripts no contêiner
RUN apk add --no-cache bash

# Estabelecemos o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copiamos o go.mod e go.sum primeiro e baixamos os módulos necessários
COPY go.mod go.sum ./
RUN go mod download

# Copiamos o restante dos arquivos fonte
COPY . .

# Construímos a aplicação
RUN go build -o /app/bin/myapp ./cmd/rest/main.go

# Exponhamos a porta 8080 para a aplicação
EXPOSE 8080

# Comando para executar a aplicação
CMD ["/app/bin/myapp"]

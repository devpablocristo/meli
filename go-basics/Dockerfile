FROM golang:1.22.1-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

WORKDIR /app/cmd/api

RUN go build -o /api

EXPOSE 8080

CMD [ "/api" ]
# Etapa 1: Construção do binário
FROM golang:1.23.3-alpine AS builder

# Instalar dependências necessárias
RUN apk add --no-cache git

# Criar diretório de trabalho
WORKDIR /app

# Copiar arquivos do projeto para dentro do container
COPY go.mod ./
COPY main.go ./

# Baixar dependências e compilar o binário
RUN go mod tidy
RUN go build -o main

# Etapa 2: Construção da imagem final
FROM alpine:latest

# Instalar dependências do runtime
RUN apk add --no-cache ca-certificates

# Configurar diretório de trabalho
WORKDIR /root/

# Copiar binário compilado da etapa anterior
COPY --from=builder /app/main .

# Expor a porta que o main usará
EXPOSE 8080

# Executar a aplicação
CMD ["./main"]

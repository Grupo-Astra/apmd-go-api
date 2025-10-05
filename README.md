# Advanced Programming & Mobile Dev - Sprint 4

## Funcionamento

É possível cadastrar sensores informando seu nome e os valores atuais lidos.
Ao cadastrar o valor atual, a API cria uma entrada no histórico daquele sensor automaticamente.

A API utiliza **PostgreSQL** como banco de dados principal.

A cada 5 segundos, por padrão, uma função executada em uma _goroutine_
atualiza o valor de todos os sensores, simulando uma leitura real.
Os novos dados também são salvos no histórico.

## Executando o Projeto (Recomendado)

A maneira mais simples e recomendada de executar o projeto é utilizando o
Docker Compose, que irá orquestrar tanto o contêiner da API quanto o
do banco de dados PostgreSQL.

### Pré-requisitos

- Docker
- Docker Compose

### Execução

1. Na raiz do projeto, execute o seguinte comando para construir as imagens
   e iniciar os contêineres:

   ```sh
   docker-compose up --build
   ```

   A API estará disponível localmente em `http://localhost:8080`.

2. Para parar e remover os contêineres e a rede criada, execute:

   ```sh
   docker-compose down
   ```

## Execução Manual (Para Desenvolvimento)

Caso prefira executar os componentes manualmente para fins de desenvolvimento.

### 1. Executar o Banco de Dados PostgreSQL

Você pode executar um contêiner do PostgreSQL separadamente:

```sh
docker run --name apmd-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=apmd_db -p 5432:5432 -d postgres:17.2-alpine
```

### 2. Executar a Aplicação Go

Com o banco de dados em execução, você pode iniciar a API:

**Exportar a variável de ambiente:**

```sh
export DATABASE_URL="postgres://postgres:password@localhost:5432/apmd_db?sslmode=disable"
```

**Executar como Script:**

```sh
go run cmd/main.go
```

**Compilar e Executar o Binário:**

```sh
# Compilar
go build -o api ./cmd

# Executar
./api
```

## Endpoints

### GET

- `/api/readings`: Retorna todos os sensores com seus valores atuais.

**Exemplo de Resposta (`/api/readings`)**

```JSON
[
    {
        "id": 1,
        "name": "Sensor de Pressão",
        "currentValue": 5.8781614244725695,
        "currentStatus": "OK",
        "historic": null
    },
    {
        "id": 2,
        "name": "Sensor de Curso (Posição)",
        "currentValue": 32.14050666057472,
        "currentStatus": "OK",
        "historic": null
    }
]

```

- `/api/readings/:id`: Retorna o sensor com id `id`, se existir,
  junto de seu histórico completo.

**Exemplo de Resposta (`/api/readings/:id`)**

```JSON
{
    "id": 1,
    "name": "Sensor de Pressão",
    "currentValue": 7.99,
    "currentStatus": "Alerta",
    "historic": [
        {
            "id": 1,
            "value": 6.2,
            "status": "OK",
            "timestamp": "2025-10-04T23:30:10.670Z",
            "sensorId": 1
        },
        {
            "id": 6,
            "value": 7.99,
            "status": "Alerta",
            "timestamp": "2025-10-04T23:30:15.706Z",
            "sensorId": 1
        }
    ]
}
```

### POST

- `/api/readings`: Cria um novo sensor, cadastrando seus valores atuais e seu histórico.

**Exemplo de Requisição**

```JSON
{
    "name": "Sensor de Temperatura",
    "currentValue": 26.0,
    "currentStatus": "OK"
}
```

**Exemplo de Resposta (Status 201 Created)**

```JSON
{
    "id": 6,
    "name": "Sensor de Temperatura",
    "currentValue": 26.0,
    "currentStatus": "OK",
    "historic": null
}
```

- `/api/database/reset`: Endpoint administrativo para resetar o banco de dados.

Esta ação:

1. Limpa **todos** os dados das tabelas de sensores e históricos.
2. Executa o "seeder" novamente para popular o banco com os dados iniciais.

**Requisição**
Não requer corpo (body).

**Exemplo de Resposta (Status 200 OK)**

```JSON
{
    "message": "Banco de dados resetado e populado com sucesso."
}
```

# Advanced Programming & Mobile Dev - Sprint 2

## Funcionamento

É possível cadastrar sensores informando seu nome e os valores atuais lidos. Ao cadastrar
o valor atual, a API cria uma entrada no histórico daquele sensor automaticamente.

A cada 5 segundos, por padrão, uma função executada em corrotina atualiza o
valor de todos os sensores, simulando uma leitura real.
Os novos dados também são salvos no histórico.

É possível executar o programa de 3 formas:

- Executando a função `main` como script
- Compilando o programa e exectando o binário
- Executar o contêiner Docker disponibilizado

## Executando o Projeto em Go

Caso tenha a linguagem Go instalada, a API pode ser executada das seguintes formas:

### Executar como Script

```sh
go run cmd/main.go
```

Se estiver na raíz do projeto, o comando executará o arquivo `main` em forma de script.

### Compilar e Executar o Binário

Também é possível compilar o programa utilizando:

```sh
go build -o api ./cmd
```

Se estiver na raíz do projeto, o comando compilará o código a
partir do diretório `cmd`, gerando um arquivo `api` na raíz.

Executar este arquivo inicializará a API.

## Buildando e Executando o Container da API

Caso não possua a Go instalado, foi disponibilizado um Dockerfile no diretório.

### Build

Monta o container com o nome `apmd-sprint2-api`:

```sh
docker build -t apmd-sprint2-api .
```

### Execução

Executa o container, expondo a API na porta `8080`:

```sh
docker run -p 8080:8080 apmd-sprint2-api
```

### Listando Dockers

Lista os contêineres em execução ou interrompidos:

```sh
docker ps -a
```

### Interrompendo Docker

Interrompe o container com id `<id>`, que pode ser verificado com `docker ps`:

```sh
docker stop <id>
```

Ou interrompa todos os contêineres:

```sh
docker stop $(docker ps -q)
```

### Removendo Docker

Remove um container interrompido de id `<id>`:

```sh
docker rm <id>
```

## Endpoints

### GET

- `/api/readings`: retorna todos os sensores com seus valores atuais

Em caso de sucesso, retorna o status 200 com um JSON do tipo:

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

- `/api/readings/:id`: retorna o sensor com id `id`, se existir, junto de seu histórico

Em caso de sucesso, retorna o status 200 com um JSON do tipo:

```JSON
{
    "id": 1,
    "name": "Sensor de Pressão",
    "currentValue": 7.996955905530946,
    "currentStatus": "Alerta",
    "historic": [
        {
            "id": 1,
            "value": 6.2,
            "status": "OK",
            "timestamp": "2025-05-27T14:07:10.670382381-03:00",
            "sensorId": 1
        },
        {
            "id": 6,
            "value": 7.893647153211203,
            "status": "Alerta",
            "timestamp": "2025-05-27T14:07:10.706692438-03:00",
            "sensorId": 1
        }
    ]
}

```

### POST

- `/api/readings`: cria um novo sensor, cadastrando seus valores atuais e seu histórico

Recebe como pedido um JSON do tipo:

```JSON
{
    "name": "Sensor de Temperatura",
    "currentValue": 26.0,
    "currentStatus": "OK"
}
```

Em caso de sucesso, retorna o status 201 com o JSON:

```JSON
{
    "id": 1,
    "name": "Sensor de Temperatura",
    "currentValue": 0,
    "currentStatus": "",
    "historic": null
}
```

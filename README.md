# Advanced Programming & Mobile Dev - Sprint 2

## Funcionamento

É possível cadastrar sensores informando seu nome e os valores atuais lidos. Ao cadastrar
o valor atual, a API cria uma entrada no histórico daquele sensor automaticamente.

A cada 5 segundos, por padrão, uma função executada em corrotina atualiza o
valor de todos os sensores, simulando uma leitura real.
Os novos dados também são salvos no histórico.

## Buildando e Executando o Container da API

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

- `/sensors`: retorna todos os sensores com seus valores atuais
- `/sensors/:id`: retorna o sensor com id `id`, se existir, junto de seu histórico

### POST

- `/sensors`: cria um novo sensor, cadastrando seus valores atuais e seu histórico

Recebe como pedido um JSON do tipo:

```JSON
{
    "name": "Sensor de Temperatura",
    "current_value": 26.0,
    "current_status": "OK"
}
```

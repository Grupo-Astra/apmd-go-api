# API de Monitoramento de Sensores - APMD Sprint 4

## 1. Introdução

Este projeto é uma API RESTful desenvolvida em Go, projetada para cadastrar sensores,
registrar suas leituras ao longo do tempo e fornecer dados através de endpoints seguros.

A aplicação simula um ambiente de IoT onde sensores são constantemente atualizados
e seus históricos são persistidos para análise.

A segurança é garantida através de autenticação baseada em tokens
JWT (JSON Web Token), e o projeto é totalmente containerizado com Docker
para facilitar a execução e a portabilidade.

### Tecnologias Utilizadas

- **Linguagem:** Go
- **Framework Web:** Gin
- **ORM:** GORM
- **Banco de Dados:** PostgreSQL
- **Autenticação:** JWT
- **Containerização:** Docker & Docker Compose

---

## 2. Como Executar o Projeto

Este projeto foi projetado para ser executado via Docker Compose,
que orquestra o serviço da API e o banco de dados PostgreSQL automaticamente.

### Pré-requisitos

- Docker
- Docker Compose

### Passos para Executar

#### 2.1 **Clone o Repositório**

```sh
git clone https://github.com/Grupo-Astra/apmd-go-api
cd apmd-go-api
```

#### 2.2 **Configure a Chave Secreta do JWT**

Antes de iniciar, você pode definir uma chave secreta customizada para a assinatura
dos tokens JWT.

Abra o arquivo `docker-compose.yml` e altere o valor da variável de ambiente `JWT_SECRET_KEY`.

```yaml
# docker-compose.yml
services:
  api:
    # ...
    environment:
      - DATABASE_URL=postgres://postgres:password@postgres:5432/apmd_db?sslmode=disable
      - JWT_SECRET_KEY: "SUA_CHAVE_SECRETA_AQUI" # <-- ALTERE ESTA LINHA
```

#### 2.3 **Inicie a Aplicação**

Execute o seguinte comando na raiz do projeto.

Ele irá construir a imagem da API, baixar a imagem do PostgreSQL e iniciar os serviços.

```sh
docker-compose up --build
```

A API estará em execução e acessível em `http://localhost:8080`.

#### 2.4 **Pare a Aplicação**

Para parar e remover os contêineres, a rede e os volumes, use o comando:

```sh
docker-compose down
```

---

## 3. Documentação da API

**Todos** os endpoints estão prefixados com `/api`.

- > **Exemplo de endpoint:** `/api/v2/readings`

As rotas de sensores são versionadas e protegidas.

### 3.1 Autenticação

Estas rotas são públicas e servem para gerenciar o acesso à API.

#### **Registrar um Novo Usuário**

- **Endpoint:** `POST /auth/register`
- **Descrição:** Cria um novo usuário no sistema. O nome de usuário deve ser único.

**Exemplo de Requisição:**

```json
{
  "username": "meu_usuario",
  "password": "minha_senha_123"
}
```

**Exemplo de Resposta (`201 Created`):**

```json
{
  "message": "Usuário criado com sucesso"
}
```

#### **Realizar Login**

- **Endpoint:** `POST /auth/login`
- **Descrição:** Autentica um usuário com base em suas credenciais e
  retorna um token JWT válido por 2 horas.

**Exemplo de Requisição:**

```json
{
  "username": "meu_usuario",
  "password": "minha_senha_123"
}
```

**Exemplo de Resposta (`200 OK`):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NjE4OTgzODZ9.abcdef123456"
}
```

### 3.2 Sensores (v2, Protegido por JWT)

Todas as requisições para estes endpoints **devem** incluir
o token JWT no cabeçalho `Authorization`.

**Exemplo de como usar o token com cURL:**

```sh
curl -X GET http://localhost:8080/api/v2/readings \
-H "Authorization: Bearer <SEU_TOKEN_JWT_AQUI>"
```

#### **Listar Todos os Sensores**

- **Endpoint:** `GET /v2/readings`
- **Descrição:** Retorna uma lista com todos os sensores
  cadastrados e seus valores atuais.

**Exemplo de Resposta (`200 OK`):**

```json
[
  {
    "ID": 1,
    "CreatedAt": "2025-10-05T15:00:00Z",
    "UpdatedAt": "2025-10-05T15:30:00Z",
    "DeletedAt": null,
    "Name": "Sensor de Pressão",
    "CurrentValue": 7.8,
    "CurrentStatus": "Alerta",
    "Historic": null
  },
  {
    "ID": 2,
    "CreatedAt": "2025-10-05T15:01:00Z",
    "UpdatedAt": "2025-10-05T15:31:00Z",
    "DeletedAt": null,
    "Name": "Contador de Ciclos",
    "CurrentValue": 1052,
    "CurrentStatus": "OK",
    "Historic": null
  }
]
```

#### **Buscar Sensor por ID**

- **Endpoint:** `GET /v2/readings/{id}`
- **Descrição:** Retorna os detalhes de um sensor específico,
  incluindo seu histórico completo de leituras.

**Exemplo de Resposta (`200 OK` para `/api/v2/readings/1`):**

```json
{
  "ID": 1,
  "CreatedAt": "2025-10-05T15:00:00Z",
  "UpdatedAt": "2025-10-05T15:30:00Z",
  "DeletedAt": null,
  "Name": "Sensor de Pressão",
  "CurrentValue": 7.8,
  "CurrentStatus": "Alerta",
  "Historic": [
    {
      "ID": 101,
      "Value": 6.2,
      "Status": "OK",
      "Timestamp": "2025-10-05T15:00:00Z",
      "SensorID": 1
    },
    {
      "ID": 152,
      "Value": 7.8,
      "Status": "Alerta",
      "Timestamp": "2025-10-05T15:30:00Z",
      "SensorID": 1
    }
  ]
}
```

#### **Criar um Novo Sensor**

- **Endpoint:** `POST /v2/readings`
- **Descrição:** Cadastra um novo sensor no sistema.
  A primeira leitura é salva automaticamente no histórico.

**Exemplo de Requisição:**

```json
{
  "name": "Sensor de Temperatura",
  "currentValue": 25.5,
  "currentStatus": "OK"
}
```

**Exemplo de Resposta (`201 Created`):**

```json
{
  "ID": 3,
  "CreatedAt": "2025-10-05T15:35:00Z",
  "UpdatedAt": "2025-10-05T15:35:00Z",
  "DeletedAt": null,
  "Name": "Sensor de Temperatura",
  "CurrentValue": 25.5,
  "CurrentStatus": "OK",
  "Historic": null
}
```

### 3.3. Administrativo

Endpoint para gerenciamento do estado da aplicação,
útil para ambientes de desenvolvimento e teste.

#### **Resetar o Banco de Dados**

- **Endpoint:** `POST /database/reset`
- **Descrição:** **Ação destrutiva!** Limpa todos os dados de todas as tabelas
  (Sensores, Históricos e Usuários) e repopula o banco com os sensores iniciais.

**Exemplo de Requisição:**
(Não requer corpo)

**Exemplo de Resposta (`200 OK`):**

```json
{
  "message": "Banco de dados resetado e populado com sucesso."
}
```

---

## 4. Arquitetura e Conceitos

Esta seção detalha a arquitetura de software do projeto e os principais padrões de
design utilizados, com o objetivo de promover um código limpo,
de fácil manutenção e testável.

### 4.1. Arquitetura do Projeto

A aplicação segue uma **Arquitetura em Camadas (Layered Architecture)**, fortemente
inspirada pelos princípios da **Clean Architecture**.

O objetivo principal é a **Separação de Preocupações (Separation of Concerns)**,
onde cada parte do sistema tem uma responsabilidade única e bem definida.

A regra fundamental é a **direção das dependências**: as camadas externas
(como a de manipulação de HTTP) dependem das camadas internas (regras
de negócio, acesso a dados), mas o inverso nunca ocorre.

Isso garante que o núcleo da aplicação seja independente de detalhes
de infraestrutura, como o framework web ou o banco de dados.

O fluxo de uma requisição HTTP protegida pode ser visualizado da seguinte forma:

```mermaid
graph TD
    subgraph "Camada Web / Infraestrutura Externa"
        A[Cliente HTTP] -- Requisição com Token --> B(Gin Router);
        B --> C{Middleware de Autenticação};
        C -- Token Inválido --> I[Resposta 401 Unauthorized];
    end

    subgraph "Camada de Aplicação / Controladores"
      C -- Token Válido --> D[Handler];
    end

    subgraph "Camada de Dados / Lógica de Negócio"
        D -- Chama --> F[Repository];
        F -- Executa Operações --> G[(GORM / PostgreSQL)];
    end

    subgraph "Fluxo de Resposta"
        G -- Retorna Dados Brutos --> F;
        F -- Mapeia para Models --> D;
        D -- Formata JSON --> H[Resposta HTTP];
        H --> A;
        I --> A;
    end
```

As camadas lógicas são:

- **Camada de Apresentação e Roteamento (Web Layer):** Composta pelos
  pacotes `main`, `routes`, `handlers` e `middleware`.
- - Sua responsabilidade é exclusivamente lidar com o protocolo HTTP.
- - O `main` inicializa e conecta tudo, o `routes` define os endpoints,
    o `middleware` intercepta e valida requisições, e os `handlers` orquestram
    o fluxo, convertendo requisições HTTP em chamadas para as camadas internas
    e formatando as respostas.

- **Camada de Acesso a Dados (Data Access Layer):** Implementada
  no pacote `repositories`.
- - Esta camada é a única que sabe como persistir e recuperar dados
    do banco de dados.
- - Ela abstrai completamente os detalhes do GORM e do
    SQL, fornecendo uma interface limpa (`SensorRepositoryInterface`,
    `UserRepositoryInterface`) para as camadas superiores.

- **Camada de Domínio e Serviços (Core Layer):** Representada pelos pacotes
  `models` e `auth`.
- - Os `models` definem as estruturas de dados centrais (as "entidades").
- - Pacotes de serviço como o `auth` contêm a lógica de negócio pura
    (ex: como gerar um hash de senha ou um token JWT), sem qualquer
    conhecimento sobre HTTP ou como os dados são salvos.

### 4.2. Conceitos Aplicados

- **Repository Pattern:** A camada de acesso a dados é abstraída através
  de repositórios. Em vez dos _handlers_ interagirem diretamente
  com o GORM, eles chamam métodos em uma interface (ex: `SensorRepositoryInterface`).

- - **Benefício:** Isso desacopla a lógica de negócio da
    implementação do banco de dados.
- - A migração do SQLite para o PostgreSQL, por exemplo, foi facilitada por essa
    abstração, pois as mudanças ficaram contidas nos repositórios, sem impactar
    o resto da aplicação.

- **Injeção de Dependência (Dependency Injection):** Os componentes não criam
  suas próprias dependências; eles as recebem prontas.
  Por exemplo, um `SensorHandler` não cria uma instância do
  `SensorRepository`, ele a recebe em seu construtor (`NewSensorHandler`).

- - **Benefício:** Aumenta a modularidade e a testabilidade.
- - Em um teste, poderíamos facilmente "injetar" um repositório
    falso (_mock_) para testar o handler de forma isolada,
    sem precisar de um banco de dados real.

- **Middleware (Gin):** Utilizamos o sistema de middleware do Gin
  para criar o `JWTAuthMiddleware`. Este é um componente que
  intercepta uma requisição antes que ela chegue ao handler final.

- - **Benefício:** É ideal para lógica de "corte transversal"
    (_cross-cutting concerns_) como autenticação,
    logging ou controle de CORS.
- - Nossa lógica de validação de token fica em um único lugar
    e é aplicada de forma limpa e declarativa a grupos
    de rotas inteiros, sem poluir os handlers de negócio.

- **Autenticação com JWT:** O sistema é _stateless_ (sem estado).
  Após o login, o servidor não armazena nenhuma informação sobre a
  sessão do usuário. O cliente recebe um token JWT e o envia a
  cada requisição para provar sua identidade.

- - **Benefício:** Facilita a escalabilidade horizontal
- - Qualquer instância da API pode validar o token sem
    precisar consultar um estado de sessão compartilhado.

- **Versionamento de API (`/v2`):** Ao introduzir uma mudança
  significativa (a obrigatoriedade de autenticação), criamos
  um novo prefixo `/v2` para as rotas protegidas.
- - **Benefício:** Permite que clientes antigos, que usavam a API
    sem autenticação, continuem funcionando (apontando para a `v1`),
    enquanto novos clientes podem usar a `v2`, mais segura.
- - É uma estratégia fundamental para evoluir APIs sem quebrar integrações
    existentes, enquanto se preparam para o processo de migração (como nosso
    frontend dos sprints anteriores).

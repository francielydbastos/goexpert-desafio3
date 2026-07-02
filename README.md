# Clean Architecture — Listagem de Orders (REST, gRPC e GraphQL)

Desafio 3 do curso Go Expert. Um único Use Case `ListOrders` é exposto simultaneamente através de três interfaces de comunicação: **REST**, **gRPC** e **GraphQL**, demonstrando o desacoplamento da Clean Architecture.

## Como executar

Requisito único: Docker e Docker Compose instalados.

```bash
docker compose up --build
```

Ao subir, o Docker Compose irá:

1. Subir o banco de dados MySQL.
2. Aguardar o banco ficar saudável (healthcheck) antes de iniciar a aplicação.
3. Aplicar automaticamente as migrações (criação da tabela `orders`).
4. Iniciar os três servidores (REST, gRPC e GraphQL).

Para parar: `docker compose down` (ou `docker compose down -v` para remover também os dados do banco).

## Portas dos serviços

| Serviço  | Porta   | Endereço                                                        |
|----------|---------|-----------------------------------------------------------------|
| REST     | `8000`  | http://localhost:8000                                           |
| GraphQL  | `8080`  | http://localhost:8080 (playground em `/`, endpoint em `/query`) |
| gRPC     | `50051` | localhost:50051                                                 |

## Use Case

- `ListOrdersUseCase` — lista todas as orders.
- `CreateOrderUseCase` — cria uma order (usado para popular o banco e testar).

## Interfaces de entrada

### REST

- `POST /order` — cria uma order.
- `GET /order` — lista as orders.

Requisições prontas no arquivo [`api.http`](./api.http).

```bash
# Criar
curl -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{"id":"a1","price":100.5,"tax":10.5}'

# Listar
curl http://localhost:8000/order
```

### GraphQL

Acesse o playground em http://localhost:8080 ou envie para `http://localhost:8080/query`.

```graphql
query {
  listOrders {
    id
    Price
    Tax
    FinalPrice
  }
}

mutation {
  createOrder(input: { id: "c3", Price: 300.0, Tax: 30.0 }) {
    id
    FinalPrice
  }
}
```

### gRPC

Service `OrderService` na porta `50051`, com reflection habilitada.

- `CreateOrder(CreateOrderRequest) returns (Order)`
- `ListOrders(Blank) returns (OrderList)`

```bash
grpcurl -plaintext -d '{"id":"d4","price":400.0,"tax":40.0}' localhost:50051 pb.OrderService/CreateOrder
grpcurl -plaintext localhost:50051 pb.OrderService/ListOrders
```

## Estrutura do projeto

```
cmd/ordersystem/        # entrypoint (wiring + wait-for-db + migracoes)
configs/                # carregamento de configuracao via env
internal/
  entity/               # Order + interface do repositorio
  usecase/              # CreateOrder, ListOrders (DTOs)
  infra/
    database/           # repositorio MySQL
    web/                # REST (handlers + webserver)
    grpc/               # proto gerado (pb) + service
    graph/              # GraphQL (gqlgen)
migrations/             # migracoes SQL
proto/                  # definicao .proto
```

## Tecnologias

- Go 1.25
- MySQL 8
- gRPC / Protocol Buffers
- GraphQL (gqlgen)
- golang-migrate (migracoes aplicadas na inicializacao, com retry para tratar race condition)
- Docker / Docker Compose

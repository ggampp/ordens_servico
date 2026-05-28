# Modelo do Banco de Dados

Banco: **PostgreSQL 16 + PostGIS 3.4**. Migrações em
`backend/internal/database/migrations`.

## Diagrama de Entidades

```
users ──────────────┐ (changed_by)
  │ employee_id      │
  ▼                  │
employees ──< employee_positions        service_order_history
  │ (responsável)                          ▲ (service_order_id)
  └──────< service_orders ─────────────────┘
```

## Tabelas

### users — autenticação
| Coluna         | Tipo        | Notas                                  |
|----------------|-------------|----------------------------------------|
| id             | bigserial   | PK                                     |
| name           | varchar     |                                        |
| email          | varchar     | único                                  |
| password_hash  | text        | bcrypt                                 |
| role           | varchar     | admin / supervisor / operator          |
| employee_id    | bigint      | FK opcional → employees                |
| active         | boolean     |                                        |
| created_at / updated_at | timestamptz |                               |

### employees — cadastro de empregados
| Coluna     | Tipo      | Notas                          |
|------------|-----------|--------------------------------|
| id         | bigserial | PK                             |
| code       | varchar   | único (Código)                 |
| name       | varchar   | Nome                           |
| email      | varchar   | E-mail                         |
| phone      | varchar   | Telefone                       |
| role       | varchar   | Cargo/Função                   |
| status     | varchar   | active / inactive              |
| deleted    | boolean   | exclusão lógica                |
| created_at / updated_at | timestamptz | Datas de cadastro/atualização |

### employee_positions — histórico de localização
| Coluna      | Tipo               | Notas                        |
|-------------|--------------------|------------------------------|
| id          | bigserial          | PK                           |
| employee_id | bigint             | FK → employees (cascade)     |
| latitude    | double precision   |                              |
| longitude   | double precision   |                              |
| geom        | geography(Point,4326) | índice GIST               |
| recorded_at | timestamptz        | Data/Hora da posição         |

### service_orders — ordens de serviço
| Coluna       | Tipo                  | Notas                                   |
|--------------|-----------------------|-----------------------------------------|
| id           | bigserial             | PK                                      |
| number       | varchar               | único (Número da Ordem)                 |
| title        | varchar               | Título                                  |
| description  | text                  | Descrição                               |
| priority     | varchar               | low / medium / high / urgent            |
| status       | varchar               | open / assigned / in_progress / completed / cancelled |
| employee_id  | bigint                | FK → employees (responsável)            |
| address      | varchar               | Endereço                                |
| latitude     | double precision      |                                         |
| longitude    | double precision      |                                         |
| geom         | geography(Point,4326) | índice GIST (filtro por região)         |
| opened_at    | timestamptz           | Data de abertura                        |
| due_at       | timestamptz           | Data prevista                           |
| completed_at | timestamptz           | Data de conclusão (automática)          |
| notes        | text                  | Observações                             |
| deleted      | boolean               | exclusão lógica                         |
| created_at / updated_at | timestamptz |                                      |

### service_order_history — histórico de status
| Coluna           | Tipo        | Notas                       |
|------------------|-------------|-----------------------------|
| id               | bigserial   | PK                          |
| service_order_id | bigint      | FK → service_orders         |
| old_status       | varchar     | status anterior (nullable)  |
| new_status       | varchar     | novo status                 |
| changed_by       | bigint      | FK → users                  |
| note             | text        | observação                  |
| changed_at       | timestamptz |                             |

## Índices relevantes

- `idx_orders_status`, `idx_orders_priority`, `idx_orders_employee` — filtros.
- `idx_orders_geom`, `idx_positions_geom` — GIST para consultas geográficas
  (filtro por região via `ST_MakeEnvelope`).
- `idx_positions_employee`, `idx_positions_recorded` — histórico/última posição.

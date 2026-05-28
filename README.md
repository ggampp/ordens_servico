# Sistema de Gestão de Ordens de Serviço com Geolocalização

Aplicação web completa para gestão de Ordens de Serviço (OS) com geolocalização,
mapa interativo e controle de equipes de campo.

Arquitetura **monolítica** e simples, pensada para um número reduzido de
usuários, baixo custo operacional e fácil manutenção.

## Stack

| Camada      | Tecnologia                                           |
|-------------|------------------------------------------------------|
| Backend     | Go 1.24 · chi · JWT · slog                           |
| Banco       | PostgreSQL 16 + PostGIS 3.4                           |
| Frontend    | React 18 · Vite · React-Leaflet (OpenStreetMap) · Tailwind |
| API         | REST · OpenAPI/Swagger                               |
| Empacotamento | Docker · Docker Compose                            |

A aplicação se conecta ao banco **exclusivamente pela URL de configuração**
(`DATABASE_URL`) do PostgreSQL/PostGIS.

## Arquitetura em camadas (backend)

```
cmd/api            -> ponto de entrada / wiring
internal/
  config           -> configuração via ambiente (DATABASE_URL, JWT, ...)
  database         -> conexão pgx + migrações embutidas (PostGIS)
  model            -> structs de domínio e DTOs
  repository       -> acesso a dados (SQL)
  service          -> regras de negócio e perfis de acesso
  handler          -> HTTP (chi) + Swagger
  middleware       -> JWT, logging estruturado, recover
  auth             -> emissão/validação de JWT
  httpx            -> tratamento centralizado de erros / validação
```

## Execução rápida (Docker Compose)

```bash
cp .env.example .env          # ajuste segredos se desejar
docker compose up --build
```

Serviços:

| Serviço   | URL                                |
|-----------|------------------------------------|
| Frontend  | http://localhost:3000              |
| Backend   | http://localhost:8080/api/v1       |
| Swagger   | http://localhost:8080/swagger      |
| Health    | http://localhost:8080/health       |

Login inicial (semeado automaticamente na primeira execução):

- **E-mail:** `admin@ordens.local`
- **Senha:** `admin123`

> As migrações (incluindo `CREATE EXTENSION postgis`) são aplicadas
> automaticamente na inicialização do backend.

## Desenvolvimento local

**Backend**

```bash
cd backend
export DATABASE_URL="postgres://ordens:ordens@localhost:5432/ordens_servico?sslmode=disable"
go run ./cmd/api
```

**Frontend**

```bash
cd frontend
npm install
npm run dev     # http://localhost:5173 (proxy /api -> :8080)
```

**Testes**

```bash
cd backend && go test ./...
```

## Perfis de acesso

| Perfil       | Permissões                                                       |
|--------------|------------------------------------------------------------------|
| Administrador| Acesso total, incluindo criação de usuários                      |
| Supervisor   | Gestão de equipes e ordens (sem criar usuários)                  |
| Operador     | Consulta e atualização **apenas** das ordens atribuídas a si     |

## Documentação

- [Manual de Implantação](docs/DEPLOYMENT.md)
- [Manual de Utilização](docs/USAGE.md)
- [Protótipos das Telas](docs/PROTOTYPES.md)
- [Modelo do Banco de Dados](docs/DATABASE.md)
- OpenAPI: `backend/internal/handler/openapi.yaml` (servido em `/swagger`)

## Endpoints principais

Autenticação: `POST /api/v1/auth/login` · `POST /api/v1/auth/register`

Empregados: `GET|POST /employees` · `GET|PUT|DELETE /employees/{id}` ·
`POST /employees/{id}/position` · `GET /employees/{id}/positions`

Ordens: `GET|POST /service-orders` · `GET|PUT|DELETE /service-orders/{id}` ·
`PATCH /service-orders/{id}/status` · `PATCH /service-orders/{id}/assign` ·
`GET /service-orders/{id}/history`

Mapa: `GET /map/employees` · `GET /map/service-orders` · `GET /map/overview`

Dashboard: `GET /dashboard`

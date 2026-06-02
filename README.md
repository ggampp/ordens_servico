# Sistema de Gestao de Ordens de Servico com Geolocalizacao

Aplicacao web para gestao de ordens de servico com geolocalizacao, mapa
interativo e controle de equipes de campo.

O projeto e um monolito: o backend Go tambem serve a SPA React compilada. Para
deploy na Forja, todo o trafego HTTP deve apontar para o servico `app` na porta
interna `8080`.

## Stack

| Camada | Tecnologia |
| --- | --- |
| Backend | Go 1.24, chi, JWT, slog |
| Banco | PostgreSQL/PostGIS gerenciado pela Forja |
| Frontend | React 18, Vite, React-Leaflet, Tailwind |
| API | REST, OpenAPI/Swagger |
| Deploy | Docker Compose na Forja |

## Deploy na Forja

O Compose segue o padrao esperado pela Forja:

- servico principal: `app`;
- porta interna: `8080`;
- bind da aplicacao: `0.0.0.0`;
- variaveis carregadas por `env_file: .env`;
- banco acessado exclusivamente por `DATABASE_URL`;
- sem `ports`, `container_name`, labels do Traefik, redes externas ou servico de banco local.

Crie ou vincule o banco pela aba Banco da Forja. Quando o banco estiver
vinculado ao app, a Forja injeta `DATABASE_URL` automaticamente.

```bash
docker compose up --build -d
```

## Variaveis

O `.env.example` contem apenas exemplos:

```env
PORT=8080
HOST=0.0.0.0
DATABASE_URL=postgres://usuario:senha@host:5432/banco?sslmode=disable
```

Outras variaveis opcionais do backend podem ser definidas no ambiente da Forja,
como `JWT_SECRET`, `JWT_EXPIRY_HOURS`, `LOG_LEVEL`, `SEED_ADMIN_EMAIL` e
`SEED_ADMIN_PASSWORD`.

## Desenvolvimento local

Para rodar localmente, configure um `.env` com uma `DATABASE_URL` valida para um
banco PostgreSQL/PostGIS acessivel pelo container ou pelo processo Go.

Backend:

```bash
cd backend
go run ./cmd/api
```

Frontend:

```bash
cd frontend
npm install
npm run dev
```

Testes:

```bash
cd backend
go test ./...
```

## Documentacao

- [Manual de Implantacao](docs/DEPLOYMENT.md)
- [Manual de Utilizacao](docs/USAGE.md)
- [Prototipos das Telas](docs/PROTOTYPES.md)
- [Modelo do Banco de Dados](docs/DATABASE.md)
- OpenAPI: `backend/internal/handler/openapi.yaml` servido em `/swagger`

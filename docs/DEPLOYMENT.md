# Manual de Implantação

## 1. Pré-requisitos

- Docker 24+ e Docker Compose v2
- (Opcional, para desenvolvimento) Go 1.24+ e Node 22+

## 2. Configuração

Copie o arquivo de exemplo e ajuste as variáveis:

```bash
cp .env.example .env
```

Variáveis relevantes:

| Variável               | Descrição                                            | Padrão                |
|------------------------|------------------------------------------------------|-----------------------|
| `DATABASE_URL`         | URL de conexão PostgreSQL/PostGIS (usada pelo backend) | montada no compose  |
| `DB_USER` / `DB_PASSWORD` / `DB_NAME` | Credenciais do Postgres               | ordens / ordens / ordens_servico |
| `JWT_SECRET`           | **Troque em produção**                               | change-me-in-production |
| `JWT_EXPIRY_HOURS`     | Validade do token                                    | 24                    |
| `SEED_ADMIN_EMAIL` / `SEED_ADMIN_PASSWORD` | Admin inicial                    | admin@ordens.local / admin123 |
| `FRONTEND_PORT` / `BACKEND_PORT` | Portas expostas no host                    | 3000 / 8080           |

> **Importante:** em produção altere `JWT_SECRET` e a senha do admin.

## 3. Subir os serviços

```bash
docker compose up --build -d
```

Isso provisiona três contêineres:

1. **db** — PostgreSQL 16 com PostGIS 3.4 (volume persistente `db_data`).
2. **backend** — API Go. Na inicialização: aguarda o banco, aplica as
   migrações (incluindo `CREATE EXTENSION postgis`) e semeia o admin.
3. **frontend** — Nginx servindo a SPA e fazendo proxy de `/api` para o backend.

## 4. Verificação

```bash
curl http://localhost:8080/health
# {"status":"ok", ...}
```

Acesse `http://localhost:3000` e faça login com o admin semeado.

## 5. Operação

```bash
docker compose logs -f backend     # logs estruturados (JSON)
docker compose ps                  # status dos serviços
docker compose down                # parar (mantém o volume de dados)
docker compose down -v             # parar e apagar dados
```

## 6. Backup e restauração do banco

```bash
# Backup
docker compose exec db pg_dump -U ordens ordens_servico > backup.sql

# Restauração
cat backup.sql | docker compose exec -T db psql -U ordens ordens_servico
```

## 7. Migrações

As migrações ficam em `backend/internal/database/migrations` e são embutidas no
binário (`go:embed`) e aplicadas idempotentemente na inicialização — não há
passo manual. Para adicionar uma nova migração, crie
`NNNNNN_descricao.up.sql` (e `.down.sql`) seguindo a numeração.

## 8. Implantação em servidor único

A arquitetura monolítica permite rodar tudo em uma única VM:

1. Instale Docker e Docker Compose.
2. Clone o repositório e configure o `.env`.
3. `docker compose up --build -d`.
4. (Recomendado) Coloque um proxy reverso (Nginx/Caddy/Traefik) com TLS à
   frente da porta do frontend.

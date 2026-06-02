# Manual de Implantacao na Forja

Este projeto esta preparado para deploy na Forja usando Docker Compose,
Traefik e banco PostgreSQL/PostGIS gerenciado pela propria Forja.

## Servico da aplicacao

O servico principal do Compose se chama `app`.

Configure a Forja para rotear o servico `app` na porta interna `8080`.
O container nao publica portas no host; ele usa apenas `expose`.

## Banco de dados

Crie ou vincule o banco PostgreSQL/PostGIS pela aba Banco da Forja.

Quando o banco for vinculado ao app, a Forja injeta a variavel
`DATABASE_URL` automaticamente no ambiente do container. A aplicacao usa
somente essa variavel para acessar o banco.

Nao configure `POSTGRES_HOST`, `POSTGRES_USER`, `POSTGRES_PASSWORD` ou
`POSTGRES_DB` no app. Tambem nao crie servico `postgres` no Compose do
projeto.

## Variaveis

O Compose carrega as variaveis geradas pela Forja com:

```yaml
env_file:
  - .env
```

O `.env.example` contem apenas exemplos:

```env
PORT=8080
HOST=0.0.0.0
DATABASE_URL=postgres://usuario:senha@host:5432/banco?sslmode=disable
```

## Deploy

A Forja deve executar:

```bash
docker compose up --build -d
```

O override da Forja injeta as redes e labels necessarias do Traefik. Por isso o
`docker-compose.yml` do projeto nao deve declarar redes externas, labels do
Traefik, `container_name` ou `ports`.

## Migrations

As migrations ficam em `backend/internal/database/migrations`, sao embutidas no
binario Go e rodam na inicializacao usando `DATABASE_URL`.

Se o banco for PostGIS, continue usando a mesma `DATABASE_URL`; a extensao deve
estar disponivel no banco criado pela Forja.

# Manual de Implantacao na Forja

Este projeto esta preparado para deploy na Forja (Hostinger VPS) com Docker
Compose, Traefik na rede `edge` e banco PostgreSQL/PostGIS gerenciado pela
propria Forja.

## Servico da aplicacao

O servico principal do Compose se chama `app`.

- Porta interna: `8080`
- Bind: `0.0.0.0` (`HOST`)
- Sem `ports`, `container_name` ou redes externas no compose base
- Health: `GET /health` ou `GET /healthz`

Na UI da Forja, defina a **porta do projeto como 8080** (o padrao e 3000). Sem
isso, o `forja-override.yml` gerado aponta para a porta errada e o Traefik
retorna 404 ou 502 mesmo com o container saudavel.

## Labels Traefik (no repositorio)

A Forja gera `forja-override.yml` com labels incompletas (porta errada, sem
TLS). O `docker-compose.yml` do projeto declara labels de TLS e porta; o
`forja-override.yml` sobrescreve a porta se o projeto estiver configurado com
valor errado na Forja — por isso a porta 8080 na UI e obrigatoria.

```yaml
labels:
  - "traefik.http.routers.forja-ordens-servico.tls=true"
  - "traefik.http.routers.forja-ordens-servico.tls.certresolver=le"
  - "traefik.http.routers.forja-ordens-servico.service=forja-ordens-servico"
  - "traefik.http.services.forja-ordens-servico.loadbalancer.server.port=8080"
```

O nome do router (`forja-ordens-servico`) deve coincidir com o slug do projeto
na Forja (`forja-<nome>`). Confira na VPS:

```bash
docker exec <container-forja> cat /data/projects/<slug>/forja-override.yml
```

Se o slug for outro, confira o valor real em `forja-override.yml` e alinhe o
nome do projeto na Forja (ou use `docker compose -p <slug>` manualmente).

Cert resolver do Traefik na VPS: `le`. Rede compartilhada: `edge`.

## Banco de dados

Crie ou vincule o banco PostgreSQL/PostGIS pela aba Banco da Forja.

Quando o banco for vinculado ao app, a Forja injeta `DATABASE_URL`
automaticamente. A aplicacao usa somente essa variavel.

Nao configure `POSTGRES_*` no app nem crie servico `postgres` no compose base.

## Variaveis

```env
PORT=8080
HOST=0.0.0.0
DATABASE_URL=postgres://usuario:senha@host:5432/banco?sslmode=disable
```

Opcionais: `JWT_SECRET`, `JWT_EXPIRY_HOURS`, `LOG_LEVEL`, `SEED_ADMIN_EMAIL`,
`SEED_ADMIN_PASSWORD`.

## Deploy

```bash
docker compose up --build -d
```

## Diagnostico na VPS (404 do Traefik)

Sintoma: containers `Running`, mas o dominio retorna `404 page not found` com
header `X-Content-Type-Options: nosniff` (resposta do Traefik, nao da app).

```bash
# Containers e portas
docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'

# Labels Traefik do app
docker inspect <container-app> --format '{{json .Config.Labels}}' | jq .

# App responde na rede edge?
docker run --rm --network edge curlimages/curl:8.5.0 -sS http://<container-app>:8080/healthz

# HTTPS pelo dominio publico
curl -sI https://<dominio>

# Logs da aplicacao
docker logs <container-app> --tail 30

# Compose mergeado pela Forja
docker exec <container-forja> cat /data/projects/<slug>/forja-override.yml
docker exec <container-forja> cat /data/projects/<slug>/docker-compose.yml
```

Interpretacao:

| Teste | OK | Problema |
| --- | --- | --- |
| curl interno `:8080/healthz` | JSON `{"status":"ok"}` | App nao sobe ou porta errada |
| curl interno falha, logs OK | — | Label `loadbalancer.server.port` errada |
| interno OK, HTTPS 404 | — | Router/rule/TLS/rede `edge` ou slug divergente |
| HTTPS 502/504 | — | App crashando ou healthcheck |

Rede e cert resolver:

```bash
docker network ls | grep -E 'edge|traefik'
docker inspect traefik --format '{{json .Args}}' | jq .
```

## Desenvolvimento local com Postgres

```bash
cp .env.example .env
docker compose -f docker-compose.yml -f docker-compose.local.yml up --build -d
curl http://localhost:8080/healthz
```

## Migrations

Em `backend/internal/database/migrations`, embutidas no binario Go e aplicadas
na inicializacao via `DATABASE_URL`.

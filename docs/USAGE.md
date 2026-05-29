# Manual de Utilização

## Acesso

1. Abra `http://localhost:8080`.
2. Faça login. O usuário inicial é `admin@ordens.local` / `admin123` (altere depois).

A navegação superior muda conforme o perfil:

- **Administrador / Supervisor:** Mapa, Dashboard, Empregados, Ordens de Serviço.
- **Operador:** Mapa e Ordens de Serviço (apenas as suas).

## Mapa Operacional

Tela inicial. Exibe sobre o OpenStreetMap:

- **Ordens de serviço** como círculos coloridos por status:
  - 🔴 Vermelho — Aberta
  - 🟠 Laranja — Atribuída
  - 🔵 Azul — Em Atendimento
  - 🟢 Verde — Concluída
  - ⚪ Cinza — Cancelada
- **Empregados** (📍) na última posição registrada.

Clique em um marcador para ver os detalhes. Em OS, há link "Ver detalhes".
Use o seletor de status para filtrar o que aparece no mapa.

## Empregados (admin/supervisor)

- **+ Novo Empregado:** preencha código, nome, e-mail, telefone, cargo, status.
- **Editar:** altera os dados (o código é imutável).
- **Excluir:** exclusão **lógica** (o registro é mantido como inativo).
- **Posição:** registra latitude/longitude manualmente ou via "Usar minha
  localização" (geolocalização do navegador). Cada registro entra no histórico.
- Filtros por status e busca por nome/código. Listagem paginada.

## Ordens de Serviço

### Criar (admin/supervisor)

"+ Nova OS": título, descrição, prioridade, responsável (opcional), endereço,
coordenadas e data prevista. O **número** é gerado automaticamente
(`OS-ANO-NNNNNN`). Se um responsável for definido na criação, a OS já nasce
como *Atribuída*.

### Listar e filtrar

Filtros por status, prioridade, responsável e intervalo de datas. Operadores
veem somente as ordens atribuídas a si.

### Detalhe da OS

- **Alterar Status:** respeita a máquina de estados:
  - Aberta → Atribuída / Em Atendimento / Cancelada
  - Atribuída → Em Atendimento / Cancelada / Aberta
  - Em Atendimento → Concluída / Cancelada
  - Concluída / Cancelada → (estados finais)
  - Ao concluir, a data de conclusão é gravada automaticamente.
- **Atribuir Responsável** (admin/supervisor).
- **Excluir** (exclusão lógica, admin/supervisor).
- **Histórico de Status:** linha do tempo de todas as mudanças, com autor,
  data e observação.

## Dashboard (admin/supervisor)

Indicadores: OS abertas, atribuídas, em atendimento, concluídas e empregados
ativos, além do gráfico de **ordens por responsável**.

## Perfis e permissões

| Ação                          | Admin | Supervisor | Operador |
|-------------------------------|:-----:|:----------:|:--------:|
| Criar usuários                |  ✅   |     ❌      |    ❌    |
| Gerenciar empregados          |  ✅   |     ✅      |    ❌    |
| Criar/atribuir/excluir OS     |  ✅   |     ✅      |    ❌    |
| Editar/atualizar status de OS |  ✅   |     ✅      | ✅ (suas) |
| Ver dashboard                 |  ✅   |     ✅      |    ❌    |
| Ver mapa                      |  ✅   |     ✅      |    ✅    |

## Documentação da API

Swagger UI interativo em `http://localhost:8080/swagger`.
Autentique-se via `POST /auth/login`, copie o `token` e use o botão
**Authorize** (Bearer) para testar os endpoints protegidos.

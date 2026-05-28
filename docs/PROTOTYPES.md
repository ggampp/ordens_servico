# Protótipos das Principais Telas

Wireframes textuais das telas implementadas no frontend (React + Tailwind).
A interface é responsiva (navegação colapsa em telas pequenas).

## 1. Login

```
        ┌──────────────────────────────┐
        │            🛠️                 │
        │     Ordens de Serviço         │
        │   Gestão com Geolocalização   │
        │                              │
        │  E-mail   [________________] │
        │  Senha    [________________] │
        │                              │
        │        [    Entrar    ]      │
        └──────────────────────────────┘
```

## 2. Mapa Operacional (tela inicial)

```
┌───────────────────────────────────────────────────────────┐
│ 🛠️ Ordens de Serviço  Mapa Dashboard Empregados Ordens  Sair│
├───────────────────────────────────────────────────────────┤
│ Mapa Operacional                    [ Todos os status ▾ ]   │
│ 🔴 Aberta 🟠 Atribuída 🔵 Em Atend. 🟢 Concluída ⚪ Canc. 📍 Emp.│
│ ┌───────────────────────────────────────────────────────┐ │
│ │                  ●(azul)        📍                      │ │
│ │       ●(vermelho)         ┌─────────────────┐          │ │
│ │                           │ OS-2026-000012  │          │ │
│ │            📍             │ Troca de poste  │          │ │
│ │                           │ Resp.: João     │          │ │
│ │     ●(verde)              │ Status: Aberta  │          │ │
│ │                           │ [Ver detalhes]  │          │ │
│ │                           └─────────────────┘          │ │
│ │             — OpenStreetMap (Leaflet) —                │ │
│ └───────────────────────────────────────────────────────┘ │
└───────────────────────────────────────────────────────────┘
```

## 3. Dashboard

```
┌───────────────────────────────────────────────────────────┐
│ Dashboard                                                   │
│ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────────┐           │
│ │  12  │ │   5  │ │   7  │ │  20  │ │    8     │           │
│ │Abertas│ │Atrib.│ │Em At.│ │Concl.│ │Emp.Ativos│           │
│ └──────┘ └──────┘ └──────┘ └──────┘ └──────────┘           │
│ ┌───────────────────────────────────────────────────────┐ │
│ │ Ordens por Responsável                                │ │
│ │  ▇▇▇▇▇▇▇  João (7)                                     │ │
│ │  ▇▇▇▇▇    Maria (5)                                    │ │
│ │  ▇▇▇      Pedro (3)                                    │ │
│ └───────────────────────────────────────────────────────┘ │
└───────────────────────────────────────────────────────────┘
```

## 4. Empregados

```
┌───────────────────────────────────────────────────────────┐
│ Empregados                              [ + Novo Empregado ]│
│ Buscar [__________]   Status [ Todos ▾ ]                    │
│ ┌───────────────────────────────────────────────────────┐ │
│ │ Código │ Nome   │ Cargo    │ Contato │ Status │ Ações  │ │
│ │ EMP001 │ João   │ Técnico  │ joão@…  │ 🟢Ativo│ Posição│ │
│ │        │        │          │         │        │ Editar │ │
│ │        │        │          │         │        │ Excluir│ │
│ │ EMP002 │ Maria  │ Eletric. │ (11)…   │ 🟢Ativo│  ...   │ │
│ └───────────────────────────────────────────────────────┘ │
│                   [Anterior]  Página 1 de 3  [Próxima]      │
└───────────────────────────────────────────────────────────┘
```

Modal "Novo/Editar": Código, Nome, E-mail, Telefone, Cargo, Status.
Modal "Posição": Latitude, Longitude, [📍 Usar minha localização].

## 5. Ordens de Serviço (lista)

```
┌───────────────────────────────────────────────────────────┐
│ Ordens de Serviço                              [ + Nova OS ]│
│ Status[▾] Prioridade[▾] Responsável[▾] De[__] Até[__]       │
│ ┌───────────────────────────────────────────────────────┐ │
│ │ Número        │ Título      │ Resp. │ Prior. │ Status   │ │
│ │ OS-2026-00012 │ Troca poste │ João  │ 🔴Alta │ 🔴Aberta │ │
│ │ OS-2026-00011 │ Reparo rede │ Maria │ 🟦Média│ 🔵Em At. │ │
│ └───────────────────────────────────────────────────────┘ │
│                   [Anterior]  Página 1 de 4  [Próxima]      │
└───────────────────────────────────────────────────────────┘
```

## 6. Detalhe da Ordem de Serviço

```
┌───────────────────────────────────────────────────────────┐
│ ← Voltar                                                    │
│ Troca de poste                          🔴 Aberta  🔴 Alta  │
│ OS-2026-000012                                              │
│ ┌─────────────────────────┬─────────────────────────────┐ │
│ │ Responsável: João       │ Endereço: Rua A, 100        │ │
│ │ Abertura: 28/05 10:00   │ Prevista: 29/05 12:00       │ │
│ │ Conclusão: —            │ Coordenadas: -23.5, -46.6   │ │
│ │ Descrição: ...                                        │ │
│ └─────────────────────────┴─────────────────────────────┘ │
│ ┌──────────────────────┐ ┌──────────────────────────────┐ │
│ │ Alterar Status       │ │ Atribuir Responsável         │ │
│ │ [ Em Atendimento ▾ ] │ │ [ Selecione… ▾ ]             │ │
│ │ [ observação ______] │ │ [   Atribuir   ]             │ │
│ │ [ Atualizar Status ] │ │ [  Excluir OS  ]             │ │
│ └──────────────────────┘ └──────────────────────────────┘ │
│ Histórico de Status                                         │
│  ● Aberta (28/05 10:00 · admin · criada)                   │
│  ● Aberta → Atribuída (28/05 10:05 · admin)                │
└───────────────────────────────────────────────────────────┘
```

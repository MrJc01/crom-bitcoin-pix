# 🗺️ Roadmap e Milestones — Crom Bitcoin Pix

> Fases de desenvolvimento com entregas concretas e critérios de aceitação.

---

## Visão Macro das Fases

```
 FASE 1         FASE 2          FASE 3          FASE 4          FASE 5
 Fundação       Lightning       Descoberta      UX/Polish       Produção
 (Wallet)       (Canais)        (Nostr)         (DX)            (Release)
 ┌──────┐      ┌──────┐        ┌──────┐        ┌──────┐        ┌──────┐
 │ M-01 │ ──▶  │ M-02 │  ──▶   │ M-03 │  ──▶   │ M-04 │  ──▶   │ M-05 │
 │      │      │      │        │      │        │      │        │      │
 │Seed  │      │LND   │        │Nostr │        │CLI   │        │Build │
 │Keys  │      │Chan  │        │LNURL │        │API   │        │Ship  │
 │Neut. │      │Pay   │        │Disc. │        │QR    │        │CI/CD │
 └──────┘      └──────┘        └──────┘        └──────┘        └──────┘
  ~2 sem        ~3 sem          ~2 sem          ~2 sem          ~1 sem
```

---

## Milestone 01: Wallet Core (Fundação)

**Objetivo:** Binário Go que gera carteira Bitcoin e consulta saldo.

### Entregas

| # | Tarefa | Detalhe |
|---|---|---|
| 1.1 | Setup do projeto Go | `go mod init`, dependências, Makefile |
| 1.2 | Geração de semente BIP-39 | 24 palavras via `crypto/rand` |
| 1.3 | Derivação de chaves BIP-32/84 | Master key → endereço SegWit |
| 1.4 | Sync Neutrino (Regtest) | Conectar ao Bitcoin Regtest local |
| 1.5 | Consulta de saldo | Verificar UTXOs via compact filters |
| 1.6 | Armazenamento local | BadgerDB para persistir estado |
| 1.7 | CLI básico | `wallet create`, `wallet balance`, `wallet restore` |

### Critério de Aceitação
```
$ crom-pay wallet create  → Gera 24 palavras + endereço bc1q...
$ crom-pay wallet balance → Mostra saldo (0 sats inicialmente)
$ crom-pay wallet restore → Restaura carteira a partir de semente
```

### Ambiente de Teste
- **Bitcoin Regtest:** Rede local para testes sem gastar BTC real
- Script `scripts/regtest.sh` sobe um bitcoind em modo regtest

---

## Milestone 02: Lightning Engine

**Objetivo:** Abrir canais Lightning e enviar/receber pagamentos instantâneos.

### Entregas

| # | Tarefa | Detalhe |
|---|---|---|
| 2.1 | LND embutido | Inicializar LND como biblioteca Go |
| 2.2 | Abertura de canal | `channels open --peer X --amount Y` |
| 2.3 | Geração de Invoice | Criar BOLT11 para recebimento |
| 2.4 | Pagamento de Invoice | Rotear pagamento via canais |
| 2.5 | Channel Monitor | Salvar estado para anti-fraude |
| 2.6 | LSP integration | Negociar inbound liquidity automático |

### Critério de Aceitação
```
# Terminal 1 (Alice)
$ crom-pay receive 5000 → Gera QR/Invoice

# Terminal 2 (Bob)
$ crom-pay pay lnbc50n1pj... → Paga em <2 segundos

# Terminal 1 (Alice)
✅ Recebido! +5000 sats
```

### Ambiente de Teste
- **2 nós LND em Regtest** (Alice + Bob)
- Canal aberto entre os dois para testes end-to-end

---

## Milestone 03: Nostr Discovery

**Objetivo:** Descobrir usuários e resolver Lightning Addresses via Nostr.

### Entregas

| # | Tarefa | Detalhe |
|---|---|---|
| 3.1 | Keypair Nostr | Derivar npub/nsec do BIP-39 |
| 3.2 | Publicar perfil | NIP-01 event nos relays |
| 3.3 | Buscar usuários | Discovery por nome/pubkey |
| 3.4 | LNURL-Pay client | Resolver `user@domain` → invoice |
| 3.5 | LNURL-Pay server | Mini HTTP para receber via Lightning Address |
| 3.6 | NIP-47 (NWC) | Controle remoto da carteira (básico) |

### Critério de Aceitação
```
$ crom-pay pay alice@crom.run 1000
🔍 Buscando alice@crom.run...
⚡ Pago! 1000 sats em 0.5s
```

---

## Milestone 04: UX e Polish

**Objetivo:** Interface CLI polida, API REST e experiência "Pix-like".

### Entregas

| # | Tarefa | Detalhe |
|---|---|---|
| 4.1 | CLI rico | Cores, spinners, tabelas formatadas |
| 4.2 | QR Code no terminal | Renderizar QR Code ASCII no CLI |
| 4.3 | API REST localhost | Endpoints para UI futura |
| 4.4 | Contatos | Salvar/gerenciar endereços frequentes |
| 4.5 | Histórico | Listar transações com filtros |
| 4.6 | Backup cifrado | Exportar/importar backup AES |

### Critério de Aceitação
- UX fluida sem erros silenciosos
- API REST documentada com Swagger/OpenAPI
- Backup restaurável em outro dispositivo

---

## Milestone 05: Release e CI/CD

**Objetivo:** Binários compilados, CI/CD e documentação final.

### Entregas

| # | Tarefa | Detalhe |
|---|---|---|
| 5.1 | Cross-compilation | Linux, macOS, Windows (amd64 + arm64) |
| 5.2 | GitHub Actions | CI: test → build → release |
| 5.3 | README.md final | Instalação, uso, contribuição |
| 5.4 | Testes E2E | Suite completa Regtest |
| 5.5 | Versão v0.1.0 | Primeiro release público |

### Critério de Aceitação
```bash
# Download e uso imediato
wget https://github.com/MrJc01/crom-bitcoin-pix/releases/download/v0.1.0/crom-pay-linux-amd64
chmod +x crom-pay-linux-amd64
./crom-pay-linux-amd64 wallet create
```

---

## Timeline Estimada

| Fase | Duração | Acumulado |
|---|---|---|
| M-01: Wallet Core | ~2 semanas | 2 semanas |
| M-02: Lightning | ~3 semanas | 5 semanas |
| M-03: Nostr Discovery | ~2 semanas | 7 semanas |
| M-04: UX/Polish | ~2 semanas | 9 semanas |
| M-05: Release | ~1 semana | 10 semanas |

> **Total estimado: ~10 semanas** para v0.1.0 funcional em Regtest/Testnet.

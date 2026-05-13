# ⚡ Crom Bitcoin Pix

> Pagamentos Bitcoin instantâneos como Pix — sem servidor, sem custódia, sem intermediários.

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/Tests-72%20passing-brightgreen)](#testes)
[![Bitcoin](https://img.shields.io/badge/Bitcoin-SegWit%20Native-orange?style=flat&logo=bitcoin)](https://github.com/bitcoin/bips/blob/master/bip-0084.mediawiki)

---

## O que é?

**Crom Bitcoin Pix** é um binário Go único que transforma Bitcoin na experiência de pagamento do Pix brasileiro:

- 💰 **Carteira HD** — BIP-39 (24 palavras) + BIP-84 (SegWit nativo `bc1q...`)
- 🔐 **Cifra militar** — AES-256-GCM + Argon2id (64MB × 3 iterações)
- ⚡ **Lightning** — Pagamentos instantâneos *(Milestone 02)*
- 🌐 **Nostr** — Identidade descentralizada *(Milestone 03)*

```bash
# Criar carteira
./crom-pay wallet create

# Ver saldo
./crom-pay wallet balance

# Exibir endereço
./crom-pay wallet address

# Restaurar de backup
./crom-pay wallet restore
```

---

## Quickstart

```bash
# Clonar
git clone https://github.com/MrJc01/crom-bitcoin-pix.git
cd crom-bitcoin-pix

# Build
make build

# Criar carteira (testnet)
./bin/crom-pay wallet create --network testnet

# Rodar testes (72 testes em 6 camadas)
make test
```

---

## Arquitetura

```
┌─────────────────────────────────────────┐
│            CLI (Cobra)                  │
├─────────────────────────────────────────┤
│         Wallet Engine                   │
│  ┌──────────┬───────────┬────────────┐  │
│  │  Seed    │ Keychain  │   Crypto   │  │
│  │ BIP-39   │  BIP-84   │ AES+Argon  │  │
│  └──────────┴───────────┴────────────┘  │
├─────────────────────────────────────────┤
│         Storage (BadgerDB)              │
└─────────────────────────────────────────┘
```

| Módulo | Função | Arquivo |
|---|---|---|
| **Seed** | Geração BIP-39 (24 palavras, 256 bits) | `internal/wallet/seed.go` |
| **Keychain** | Derivação BIP-84 (`m/84'/0'/0'/0/i`) | `internal/wallet/keychain.go` |
| **Crypto** | AES-256-GCM + Argon2id | `internal/wallet/crypto.go` |
| **Wallet** | Orquestração create/open/restore | `internal/wallet/wallet.go` |
| **Storage** | BadgerDB wrapper | `internal/storage/db.go` |

---

## Segurança

| Aspecto | Implementação |
|---|---|
| Entropia | `crypto/rand` (CSPRNG do OS) — Shannon 7.98/8.0 |
| KDF | Argon2id (64MB RAM, 3 iter, 4 threads) |
| Cifra | AES-256-GCM (AEAD: confidencialidade + integridade) |
| Endereços | P2WPKH SegWit nativo (Bech32) |
| Permissões | Diretório wallet.db com mode `0700` |
| Memória | `zeroBytes()` + `Keychain.Zero()` + `runtime.KeepAlive` |
| Brute-force | ~1 tentativa/segundo |
| BIP-84 | ✅ Vetor oficial verificado: `bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu` |

---

## Testes

**72 testes** em **6 camadas**, todos passando:

```bash
# Rodar tudo (gera relatório .md em tests/reports/)
./tests/run_all.sh
```

| Camada | Testes | O que cobre |
|---|---|---|
| Unit (Wallet) | 24 | Seed, keychain, crypto, BIP-84 vector |
| Unit (Storage) | 12 | CRUD, persistência, permissões, prefixos |
| Integration | 6 | Create→Open→Restore, redes, persistência |
| Security | 15 | Entropia, timing, tamper, key exposure |
| Stress | 4 | 20 wallets concorrentes, 1000 derivações |
| E2E (CLI) | 11 | Binário real: help, create, balance, address |

---

## Documentação

📚 Documentação completa disponível em [`documentacao/`](documentacao/):

| Documento | Conteúdo |
|---|---|
| [Início](documentacao/01-inicio.md) | Visão geral e quickstart |
| [Arquitetura](documentacao/02-arquitetura.md) | Design técnico e módulos |
| [Criptografia](documentacao/03-criptografia.md) | Primitivas e segurança |
| [Comandos CLI](documentacao/04-comandos.md) | Referência de todos os comandos |
| [Testes](documentacao/05-testes.md) | Estrutura e como rodar |
| [Roadmap](documentacao/06-roadmap.md) | Milestones e próximos passos |
| [Glossário](documentacao/07-glossario.md) | Termos técnicos |
| [Contribuição](documentacao/08-contribuicao.md) | Como contribuir |

---

## Estrutura do Projeto

```
crom-bitcoin-pix/
├── cmd/crom-pay/          # CLI (main, root, wallet commands)
├── internal/
│   ├── wallet/            # Seed, keychain, crypto, wallet
│   └── storage/           # BadgerDB wrapper
├── tests/
│   ├── integration/       # Fluxos completos
│   ├── security/          # Pentest automatizado
│   ├── stress/            # Carga e concorrência
│   ├── e2e/               # Testes do binário
│   ├── reports/           # Relatórios gerados
│   └── run_all.sh         # Runner principal
├── documentacao/           # Documentação completa (PT-BR)
├── docs/                  # Specs técnicas originais
├── Makefile               # Build, test, clean
└── LICENSE                # MIT
```

---

## Roadmap

| Milestone | Status | Descrição |
|---|---|---|
| **01 — Wallet Core** | ✅ Completo | BIP-39/84, AES-GCM, CLI, 72 testes |
| **02 — Lightning** | 🔜 Próximo | LND embutido, Neutrino, invoices |
| **03 — Nostr Identity** | ⬜ Planejado | Descoberta P2P, endereços legíveis |
| **04 — Pix UX** | ⬜ Planejado | QR Code, TUI, pagamento em 1 toque |

---

## Licença

[MIT](LICENSE) — © 2026 MrJc01 / Crom Project

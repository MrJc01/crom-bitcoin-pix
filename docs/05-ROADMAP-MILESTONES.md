# 🗺️ Roadmap e Milestones — Crom Bitcoin Pix

> Fases de desenvolvimento com entregas concretas e critérios de aceitação.

---

## Visão Macro das Fases

```
 FASE 1         FASE 2          FASE 3          FASE 4          FASE 5
 Fundação       Lightning       Descoberta      UX/Polish       CI/CD
 (Wallet)       (Canais)        (Nostr)         (DX)            (Release)
 ┌──────┐      ┌──────┐        ┌──────┐        ┌──────┐        ┌──────┐
 │ M-01 │ ──▶  │ M-02 │  ──▶   │ M-03 │  ──▶   │ M-04 │  ──▶   │ M-05 │
 │  ✅  │      │  ✅  │        │  ✅  │        │  ✅  │        │  ✅  │
 │Seed  │      │LND   │        │Nostr │        │CLI   │        │Build │
 │Keys  │      │Chain │        │NIP05 │        │QR/TUI│        │Ship  │
 │Crypt │      │Pay   │        │Zaps  │        │Contac│        │CI/CD │
 └──────┘      └──────┘        └──────┘        └──────┘        └──────┘
  ✅ Done       ✅ Done         ✅ Done         ✅ Done         ✅ Done
```

---

## ✅ Milestone 01: Wallet Core (Completo)

**Objetivo:** Binário Go que gera carteira Bitcoin com criptografia auditada.

| # | Tarefa | Status |
|---|---|---|
| 1.1 | Setup do projeto Go + Makefile | ✅ |
| 1.2 | Geração de semente BIP-39 (24 palavras, crypto/rand) | ✅ |
| 1.3 | Derivação BIP-32/84 → endereço SegWit | ✅ |
| 1.4 | AES-256-GCM + Argon2id (64MB/3 iter) | ✅ |
| 1.5 | BadgerDB storage com prefix isolation | ✅ |
| 1.6 | CLI: `wallet create/balance/address/restore` | ✅ |
| 1.7 | Auditoria de segurança (6 vulnerabilidades corrigidas) | ✅ |

---

## ✅ Milestone 02: Lightning Engine (Completo)

**Objetivo:** Consulta de saldo real e infraestrutura Lightning.

| # | Tarefa | Status |
|---|---|---|
| 2.1 | Esplora API client (mempool.space) | ✅ |
| 2.2 | Saldo real (confirmed + unconfirmed) | ✅ |
| 2.3 | TX Builder P2WPKH (UTXO selection, witness signing, RBF) | ✅ |
| 2.4 | LND REST Client (macaroon auth, TLS) | ✅ |
| 2.5 | Invoice BOLT-11 (criar, decodificar, pagar) | ✅ |
| 2.6 | Channel management (open, close, list) | ✅ |
| 2.7 | CLI: `lightning info/invoice/channels/setup` | ✅ |
| 2.8 | CLI: `pay` (BTC/Lightning/NIP-05) + `receive` (QR) | ✅ |

---

## ✅ Milestone 03: Nostr Identity (Completo)

**Objetivo:** Identidade descentralizada para descoberta P2P.

| # | Tarefa | Status |
|---|---|---|
| 3.1 | Derivação NIP-06 (m/44'/1237'/0'/0/0) | ✅ |
| 3.2 | Relay pool com fallback (5 relays default) | ✅ |
| 3.3 | Publicar perfil (kind 0) e notas (kind 1) | ✅ |
| 3.4 | Verificação NIP-05 (user@domain) | ✅ |
| 3.5 | Zap request (kind 9734) | ✅ |
| 3.6 | npub/nsec em bech32 | ✅ |
| 3.7 | CLI: `nostr identity/publish/relays/verify` | ✅ |

---

## ✅ Milestone 04: Pix UX (Completo)

**Objetivo:** Experiência de pagamento Pix-like no terminal.

| # | Tarefa | Status |
|---|---|---|
| 4.1 | QR Code ASCII (Unicode blocks, normal + invertido) | ✅ |
| 4.2 | URI BIP-21 (`bitcoin:bc1q...?amount=0.001`) | ✅ |
| 4.3 | TUI bubbletea + lipgloss (5 telas, tema Bitcoin) | ✅ |
| 4.4 | Contatos CRUD (BadgerDB + resolução automática) | ✅ |
| 4.5 | FormatSats (sats, mBTC, BTC) | ✅ |
| 4.6 | CLI: `contacts add/list/show/remove` | ✅ |
| 4.7 | CLI: `tui` (dashboard interativo) | ✅ |

---

## ✅ Milestone 05: CI/CD e Release (Completo)

**Objetivo:** Automação, distribuição e qualidade.

| # | Tarefa | Status |
|---|---|---|
| 5.1 | GitHub Actions CI (test + vet + build) | ✅ |
| 5.2 | Build matrix: linux/darwin amd64+arm64, windows | ✅ |
| 5.3 | Release automático com SHA256 (via tags) | ✅ |
| 5.4 | 79+ testes em 6 camadas (100% pass) | ✅ |
| 5.5 | README completo com 30+ comandos documentados | ✅ |

```bash
# Download e uso imediato (após primeira release)
wget https://github.com/MrJc01/crom-bitcoin-pix/releases/download/v0.1.0/crom-pay-linux-amd64.tar.gz
tar xzf crom-pay-linux-amd64.tar.gz
./crom-pay-linux-amd64 wallet create
```

---

## ⬜ Milestone 06: Multi-plataforma (Futuro)

| # | Tarefa | Descrição |
|---|---|---|
| 6.1 | Mobile (Gomobile) | Android/iOS via bindings |
| 6.2 | Web (WASM) | Interface web |
| 6.3 | Desktop GUI | GUI nativa |
| 6.4 | Docker | Container para servidor |

---

## Timeline

| Fase | Status |
|---|---|
| M-01: Wallet Core | ✅ Completo |
| M-02: Lightning Engine | ✅ Completo |
| M-03: Nostr Identity | ✅ Completo |
| M-04: Pix UX | ✅ Completo |
| M-05: CI/CD | ✅ Completo |
| M-06: Multi-plataforma | ⬜ Futuro |
